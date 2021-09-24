// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"

	"go.temporal.io/server/common/dynamicconfig"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/persistence"
	"go.temporal.io/server/common/persistence/visibility/manager"
	"go.temporal.io/server/common/persistence/visibility/store"
	"go.temporal.io/server/common/persistence/visibility/store/elasticsearch/client"
	"go.temporal.io/server/common/persistence/visibility/store/elasticsearch/query"
	"go.temporal.io/server/common/searchattribute"
)

const (
	persistenceName = "elasticsearch"

	delimiter                    = "~"
	pointInTimeKeepAliveInterval = "1m"
	scrollKeepAliveInterval      = "1m"
)

type (
	visibilityStore struct {
		esClient                 client.Client
		index                    string
		searchAttributesProvider searchattribute.Provider
		searchAttributesMapper   searchattribute.Mapper
		processor                Processor
		processorAckTimeout      dynamicconfig.DurationPropertyFn
		metricsClient            metrics.Client
	}

	visibilityPageToken struct {
		SearchAfter []interface{}

		// For ScanWorkflowExecutions API.
		// For ES<7.10.0 and "oss" flavor.
		ScrollID string
		// For ES>=7.10.0 and "default" flavor.
		PointInTimeID string
	}
)

var _ store.VisibilityStore = (*visibilityStore)(nil)

var (
	errUnexpectedJSONFieldType = errors.New("unexpected JSON field type")
)

// NewVisibilityStore create a visibility store connecting to ElasticSearch
func NewVisibilityStore(
	esClient client.Client,
	index string,
	searchAttributesProvider searchattribute.Provider,
	searchAttributesMapper searchattribute.Mapper,
	processor Processor,
	processorAckTimeout dynamicconfig.DurationPropertyFn,
	metricsClient metrics.Client,
) *visibilityStore {

	return &visibilityStore{
		esClient:                 esClient,
		index:                    index,
		searchAttributesProvider: searchAttributesProvider,
		searchAttributesMapper:   searchAttributesMapper,
		processor:                processor,
		processorAckTimeout:      processorAckTimeout,
		metricsClient:            metricsClient,
	}
}

func (s *visibilityStore) Close() {
	// TODO (alex): visibilityStore shouldn't Stop processor. Processor should be stopped where it is created.
	if s.processor != nil {
		s.processor.Stop()
	}
}

func (s *visibilityStore) GetName() string {
	return persistenceName
}

func (s *visibilityStore) RecordWorkflowExecutionStarted(request *store.InternalRecordWorkflowExecutionStartedRequest) error {
	visibilityTaskKey := getVisibilityTaskKey(request.ShardID, request.TaskID)
	doc, err := s.generateESDoc(request.InternalVisibilityRequestBase, visibilityTaskKey)
	if err != nil {
		return err
	}

	return s.addBulkIndexRequestAndWait(request.InternalVisibilityRequestBase, doc, visibilityTaskKey)
}

func (s *visibilityStore) RecordWorkflowExecutionClosed(request *store.InternalRecordWorkflowExecutionClosedRequest) error {
	visibilityTaskKey := getVisibilityTaskKey(request.ShardID, request.TaskID)
	doc, err := s.generateESDoc(request.InternalVisibilityRequestBase, visibilityTaskKey)
	if err != nil {
		return err
	}

	doc[searchattribute.CloseTime] = request.CloseTime
	doc[searchattribute.ExecutionDuration] = request.CloseTime.Sub(request.ExecutionTime).Nanoseconds()
	doc[searchattribute.HistoryLength] = request.HistoryLength
	doc[searchattribute.StateTransitionCount] = request.StateTransitionCount

	return s.addBulkIndexRequestAndWait(request.InternalVisibilityRequestBase, doc, visibilityTaskKey)
}

func (s *visibilityStore) UpsertWorkflowExecution(request *store.InternalUpsertWorkflowExecutionRequest) error {
	visibilityTaskKey := getVisibilityTaskKey(request.ShardID, request.TaskID)
	doc, err := s.generateESDoc(request.InternalVisibilityRequestBase, visibilityTaskKey)
	if err != nil {
		return err
	}

	return s.addBulkIndexRequestAndWait(request.InternalVisibilityRequestBase, doc, visibilityTaskKey)
}

func (s *visibilityStore) DeleteWorkflowExecution(request *manager.VisibilityDeleteWorkflowExecutionRequest) error {
	docID := getDocID(request.WorkflowID, request.RunID)

	bulkDeleteRequest := &client.BulkableRequest{
		Index:       s.index,
		ID:          docID,
		Version:     request.TaskID,
		RequestType: client.BulkableRequestTypeDelete,
	}

	return s.addBulkRequestAndWait(bulkDeleteRequest, docID)
}

func getDocID(workflowID string, runID string) string {
	return fmt.Sprintf("%s%s%s", workflowID, delimiter, runID)
}

func getVisibilityTaskKey(shardID int32, taskID int64) string {
	return fmt.Sprintf("%d%s%d", shardID, delimiter, taskID)
}

func (s *visibilityStore) addBulkIndexRequestAndWait(
	request *store.InternalVisibilityRequestBase,
	esDoc map[string]interface{},
	visibilityTaskKey string,
) error {
	bulkIndexRequest := &client.BulkableRequest{
		Index:       s.index,
		ID:          getDocID(request.WorkflowID, request.RunID),
		Version:     request.TaskID,
		RequestType: client.BulkableRequestTypeIndex,
		Doc:         esDoc,
	}

	return s.addBulkRequestAndWait(bulkIndexRequest, visibilityTaskKey)
}

func (s *visibilityStore) addBulkRequestAndWait(bulkRequest *client.BulkableRequest, visibilityTaskKey string) error {
	s.checkProcessor()

	ackCh := s.processor.Add(bulkRequest, visibilityTaskKey)
	ackTimeoutTimer := time.NewTimer(s.processorAckTimeout())
	defer ackTimeoutTimer.Stop()

	select {
	case ack := <-ackCh:
		if !ack {
			return newVisibilityTaskNAckError(visibilityTaskKey)
		}
		return nil
	case <-ackTimeoutTimer.C:
		return newVisibilityTaskAckTimeoutError(visibilityTaskKey, s.processorAckTimeout())
	}
}

func (s *visibilityStore) checkProcessor() {
	if s.processor == nil {
		// must be bug, check history setup
		panic("elastic search processor is nil")
	}
	if s.processorAckTimeout == nil {
		// must be bug, check history setup
		panic("config.ESProcessorAckTimeout is nil")
	}
}

func (s *visibilityStore) ListOpenWorkflowExecutions(request *manager.ListWorkflowExecutionsRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request, boolQuery, true)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListOpenWorkflowExecutions failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.StartTime.Before(request.EarliestStartTime) && !rec.StartTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListClosedWorkflowExecutions(request *manager.ListWorkflowExecutionsRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		MustNot(elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request, boolQuery, false)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListClosedWorkflowExecutions failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.CloseTime.Before(request.EarliestStartTime) && !rec.CloseTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListOpenWorkflowExecutionsByType(request *manager.ListWorkflowExecutionsByTypeRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(
			elastic.NewTermQuery(searchattribute.WorkflowType, request.WorkflowTypeName),
			elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request.ListWorkflowExecutionsRequest, boolQuery, true)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListOpenWorkflowExecutionsByType failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.StartTime.Before(request.EarliestStartTime) && !rec.StartTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListClosedWorkflowExecutionsByType(request *manager.ListWorkflowExecutionsByTypeRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(elastic.NewTermQuery(searchattribute.WorkflowType, request.WorkflowTypeName)).
		MustNot(elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request.ListWorkflowExecutionsRequest, boolQuery, false)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListClosedWorkflowExecutionsByType failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.CloseTime.Before(request.EarliestStartTime) && !rec.CloseTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListOpenWorkflowExecutionsByWorkflowID(request *manager.ListWorkflowExecutionsByWorkflowIDRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(
			elastic.NewTermQuery(searchattribute.WorkflowID, request.WorkflowID),
			elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request.ListWorkflowExecutionsRequest, boolQuery, true)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListOpenWorkflowExecutionsByWorkflowID failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.StartTime.Before(request.EarliestStartTime) && !rec.StartTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListClosedWorkflowExecutionsByWorkflowID(request *manager.ListWorkflowExecutionsByWorkflowIDRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(elastic.NewTermQuery(searchattribute.WorkflowID, request.WorkflowID)).
		MustNot(elastic.NewTermQuery(searchattribute.ExecutionStatus, enumspb.WORKFLOW_EXECUTION_STATUS_RUNNING.String()))

	p, err := s.buildSearchParameters(request.ListWorkflowExecutionsRequest, boolQuery, false)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListClosedWorkflowExecutionsByWorkflowID failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.CloseTime.Before(request.EarliestStartTime) && !rec.CloseTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListClosedWorkflowExecutionsByStatus(request *manager.ListClosedWorkflowExecutionsByStatusRequest) (*store.InternalListWorkflowExecutionsResponse, error) {

	boolQuery := elastic.NewBoolQuery().
		Filter(elastic.NewTermQuery(searchattribute.ExecutionStatus, request.Status.String()))

	p, err := s.buildSearchParameters(request.ListWorkflowExecutionsRequest, boolQuery, false)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListClosedWorkflowExecutionsByStatus failed. Error: %s", detailedErrorMessage(err)))
	}

	isRecordValid := func(rec *store.InternalWorkflowExecutionInfo) bool {
		return !rec.CloseTime.Before(request.EarliestStartTime) && !rec.CloseTime.After(request.LatestStartTime)
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, isRecordValid)
}

func (s *visibilityStore) ListWorkflowExecutions(request *manager.ListWorkflowExecutionsRequestV2) (*store.InternalListWorkflowExecutionsResponse, error) {
	p, err := s.buildSearchParametersV2(request)
	if err != nil {
		return nil, err
	}

	token, err := s.deserializePageToken(request.NextPageToken)
	if err != nil {
		return nil, err
	}

	if token != nil && len(token.SearchAfter) > 0 {
		p.SearchAfter = token.SearchAfter
	}

	ctx := context.Background()
	searchResult, err := s.esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ListWorkflowExecutions failed. Error: %s", detailedErrorMessage(err)))
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, nil)
}

func (s *visibilityStore) ScanWorkflowExecutions(request *manager.ListWorkflowExecutionsRequestV2) (*store.InternalListWorkflowExecutionsResponse, error) {
	ctx := context.Background()

	if esClientV7, isV7 := s.esClient.(client.ClientV7); isV7 {
		// Elasticsearch 7.10+ can use "point in time" (PIT) instead of scroll to scan over all workflows without skipping or duplicating them.
		// https://www.elastic.co/guide/en/elasticsearch/reference/7.10/point-in-time-api.html
		if esClientV7.IsPointInTimeSupported(ctx) {
			return s.scanWorkflowExecutionsWithPit(ctx, request, esClientV7)
		}
	}

	return s.scanWorkflowExecutionsWithScroll(ctx, request)
}

func (s *visibilityStore) scanWorkflowExecutionsWithPit(ctx context.Context, request *manager.ListWorkflowExecutionsRequestV2, esClient client.ClientV7) (*store.InternalListWorkflowExecutionsResponse, error) {
	p, err := s.buildSearchParametersV2(request)
	if err != nil {
		return nil, err
	}

	// First call doesn't have token with PointInTimeID.
	if len(request.NextPageToken) == 0 {
		pitID, err := esClient.OpenPointInTime(ctx, s.index, pointInTimeKeepAliveInterval)
		if err != nil {
			return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to create point in time: %s", detailedErrorMessage(err)))
		}
		p.PointInTime = elastic.NewPointInTimeWithKeepAlive(pitID, pointInTimeKeepAliveInterval)
	} else {
		token, err := s.deserializePageToken(request.NextPageToken)
		if err != nil {
			return nil, err
		}
		if token.PointInTimeID == "" {
			return nil, serviceerror.NewInvalidArgument("pointInTimeId must present in pagination token")
		}
		p.SearchAfter = token.SearchAfter
		p.PointInTime = elastic.NewPointInTimeWithKeepAlive(token.PointInTimeID, pointInTimeKeepAliveInterval)
	}

	searchResult, err := esClient.Search(ctx, p)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ScanWorkflowExecutions failed. Error: %s", detailedErrorMessage(err)))
	}

	// Empty hits list indicate that this is a last page.
	if searchResult.Hits != nil && len(searchResult.Hits.Hits) < request.PageSize {
		_, err = esClient.ClosePointInTime(ctx, searchResult.PitId)
		if err != nil {
			return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to close point in time: %s", detailedErrorMessage(err)))
		}
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, nil)
}

func (s *visibilityStore) scanWorkflowExecutionsWithScroll(ctx context.Context, request *manager.ListWorkflowExecutionsRequestV2) (*store.InternalListWorkflowExecutionsResponse, error) {
	var (
		searchResult *elastic.SearchResult
		scrollErr    error
	)

	// First call doesn't have token with ScrollID.
	if len(request.NextPageToken) == 0 {
		// First page.
		p, err := s.buildSearchParametersV2(request)
		if err != nil {
			return nil, err
		}
		searchResult, scrollErr = s.esClient.OpenScroll(ctx, p, scrollKeepAliveInterval)
	} else {
		token, err := s.deserializePageToken(request.NextPageToken)
		if err != nil {
			return nil, err
		}
		if token.ScrollID == "" {
			return nil, serviceerror.NewInvalidArgument("scrollId must present in pagination token")
		}
		searchResult, scrollErr = s.esClient.Scroll(ctx, token.ScrollID, scrollKeepAliveInterval)
	}

	if scrollErr != nil && scrollErr != io.EOF {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("ScanWorkflowExecutions failed. Error: %s", detailedErrorMessage(scrollErr)))
	}

	// Both io.IOF and empty hits list indicate that this is a last page.
	if (searchResult.Hits != nil && len(searchResult.Hits.Hits) < request.PageSize) ||
		scrollErr == io.EOF {
		err := s.esClient.CloseScroll(ctx, searchResult.ScrollId)
		if err != nil {
			return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to close scroll: %s", detailedErrorMessage(err)))
		}
	}

	return s.getListWorkflowExecutionsResponse(searchResult, request.Namespace, request.PageSize, nil)
}

func (s *visibilityStore) CountWorkflowExecutions(request *manager.CountWorkflowExecutionsRequest) (*manager.CountWorkflowExecutionsResponse, error) {
	boolQuery, _, err := s.convertQuery(request.Namespace, request.NamespaceID, request.Query)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	count, err := s.esClient.Count(ctx, s.index, boolQuery)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("CountWorkflowExecutions failed. Error: %s", detailedErrorMessage(err)))
	}

	response := &manager.CountWorkflowExecutionsResponse{Count: count}
	return response, nil
}

func (s *visibilityStore) buildSearchParameters(
	request *manager.ListWorkflowExecutionsRequest,
	boolQuery *elastic.BoolQuery,
	overStartTime bool,
) (*client.SearchParameters, error) {

	token, err := s.deserializePageToken(request.NextPageToken)
	if err != nil {
		return nil, err
	}

	boolQuery.Filter(elastic.NewTermQuery(searchattribute.NamespaceID, request.NamespaceID))

	if !request.EarliestStartTime.IsZero() || !request.LatestStartTime.IsZero() {
		var rangeQuery *elastic.RangeQuery
		if overStartTime {
			rangeQuery = elastic.NewRangeQuery(searchattribute.StartTime)
		} else {
			rangeQuery = elastic.NewRangeQuery(searchattribute.CloseTime)
		}

		if !request.EarliestStartTime.IsZero() {
			rangeQuery = rangeQuery.Gte(request.EarliestStartTime)
		}

		if !request.LatestStartTime.IsZero() {
			rangeQuery = rangeQuery.Lte(request.LatestStartTime)
		}
		boolQuery.Filter(rangeQuery)
	}

	params := &client.SearchParameters{
		Index:    s.index,
		Query:    boolQuery,
		PageSize: request.PageSize,
	}

	if token != nil && len(token.SearchAfter) > 0 {
		params.SearchAfter = token.SearchAfter
	}

	if overStartTime {
		params.Sorter = append(params.Sorter, elastic.NewFieldSort(searchattribute.StartTime).Desc())
	} else {
		params.Sorter = append(params.Sorter, elastic.NewFieldSort(searchattribute.CloseTime).Desc())
	}

	// RunID is explicit tiebreaker.
	params.Sorter = append(params.Sorter, elastic.NewFieldSort(searchattribute.RunID).Desc())

	return params, nil
}

func (s *visibilityStore) buildSearchParametersV2(
	request *manager.ListWorkflowExecutionsRequestV2,
) (*client.SearchParameters, error) {

	boolQuery, fieldSorts, err := s.convertQuery(request.Namespace, request.NamespaceID, request.Query)
	if err != nil {
		return nil, err
	}

	params := &client.SearchParameters{
		Index:    s.index,
		Query:    boolQuery,
		PageSize: request.PageSize,
		Sorter:   s.setDefaultFieldSort(fieldSorts),
	}

	return params, nil
}
func (s *visibilityStore) convertQuery(namespace string, namespaceID string, requestQueryStr string) (*elastic.BoolQuery, []*elastic.FieldSort, error) {
	saTypeMap, err := s.searchAttributesProvider.GetSearchAttributes(s.index, false)
	if err != nil {
		return nil, nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to read search attribute types: %v", err))
	}
	queryConverter := query.NewConverter(
		newNameInterceptor(namespace, s.index, saTypeMap, s.searchAttributesMapper),
		newValuesInterceptor(),
	)
	requestQuery, fieldSorts, err := queryConverter.ConvertWhereOrderBy(requestQueryStr)
	if err != nil {
		// Convert query.ConverterError to serviceerror.InvalidArgument and pass through all other errors (which should be only mapper errors).
		var converterErr *query.ConverterError
		if errors.As(err, &converterErr) {
			return nil, nil, serviceerror.NewInvalidArgument(fmt.Sprintf("unable to parse query: %v", converterErr))
		}
		return nil, nil, err
	}

	// Create new bool query because request query might have only "should" (="or") queries.
	namespaceFilterQuery := elastic.NewBoolQuery().Filter(elastic.NewTermQuery(searchattribute.NamespaceID, namespaceID))
	if requestQuery != nil {
		namespaceFilterQuery.Filter(requestQuery)
	}

	return namespaceFilterQuery, fieldSorts, nil
}

func (s *visibilityStore) setDefaultFieldSort(fieldSorts []*elastic.FieldSort) []elastic.Sorter {
	if len(fieldSorts) == 0 {
		// Set default sorting by StartTime desc and RunID as tiebreaker.
		return []elastic.Sorter{
			elastic.NewFieldSort(searchattribute.StartTime).Desc(),
			elastic.NewFieldSort(searchattribute.RunID).Desc(),
		}
	}

	res := make([]elastic.Sorter, len(fieldSorts)+1)
	for i, fs := range fieldSorts {
		res[i] = fs
	}
	// RunID is explicit tiebreaker.
	res[len(res)-1] = elastic.NewFieldSort(searchattribute.RunID).Desc()

	return res
}

func (s *visibilityStore) getListWorkflowExecutionsResponse(
	searchResult *elastic.SearchResult,
	namespace string,
	pageSize int,
	isRecordValid func(rec *store.InternalWorkflowExecutionInfo) bool,
) (*store.InternalListWorkflowExecutionsResponse, error) {

	if searchResult.Hits == nil || len(searchResult.Hits.Hits) == 0 {
		return &store.InternalListWorkflowExecutionsResponse{}, nil
	}

	typeMap, err := s.searchAttributesProvider.GetSearchAttributes(s.index, false)
	if err != nil {
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to read search attribute types: %v", err))
	}

	response := &store.InternalListWorkflowExecutionsResponse{
		Executions: make([]*store.InternalWorkflowExecutionInfo, 0, len(searchResult.Hits.Hits)),
	}
	var lastHitSort []interface{}
	for _, hit := range searchResult.Hits.Hits {
		workflowExecutionInfo, err := s.parseESDoc(hit, typeMap, namespace)
		if err != nil {
			return nil, err
		}
		// ES6 uses "date" data type not "date_nanos". It truncates dates using milliseconds and might return extra rows.
		// For example: 2021-06-12T00:21:43.159739259Z fits 2021-06-12T00:21:43.158Z...2021-06-12T00:21:43.159Z range lte/gte query.
		// Therefore, these records need to be filtered out on the client side to support nanos precision.
		// After ES6 deprecation isRecordValid can be removed.
		if isRecordValid == nil || isRecordValid(workflowExecutionInfo) {
			response.Executions = append(response.Executions, workflowExecutionInfo)
			lastHitSort = hit.Sort
		}
	}

	if len(searchResult.Hits.Hits) == pageSize && lastHitSort != nil { // this means the response is not the last page
		response.NextPageToken, err = s.serializePageToken(&visibilityPageToken{
			SearchAfter:   lastHitSort,
			PointInTimeID: searchResult.PitId,
			ScrollID:      searchResult.ScrollId,
		})
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (s *visibilityStore) deserializePageToken(data []byte) (*visibilityPageToken, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var token *visibilityPageToken
	dec := json.NewDecoder(bytes.NewReader(data))
	// UseNumber will not lose precision on big int64.
	dec.UseNumber()
	err := dec.Decode(&token)
	if err != nil {
		return nil, serviceerror.NewInvalidArgument(fmt.Sprintf("unable to deserialize page token: %v", err))
	}
	return token, nil
}

func (s *visibilityStore) serializePageToken(token *visibilityPageToken) ([]byte, error) {
	if token == nil {
		return nil, nil
	}

	data, err := json.Marshal(token)
	if err != nil {
		return nil, serviceerror.NewInvalidArgument(fmt.Sprintf("unable to serialize page token: %v", err))
	}
	return data, nil
}

func (s *visibilityStore) generateESDoc(request *store.InternalVisibilityRequestBase, visibilityTaskKey string) (map[string]interface{}, error) {
	doc := map[string]interface{}{
		searchattribute.VisibilityTaskKey: visibilityTaskKey,
		searchattribute.NamespaceID:       request.NamespaceID,
		searchattribute.WorkflowID:        request.WorkflowID,
		searchattribute.RunID:             request.RunID,
		searchattribute.WorkflowType:      request.WorkflowTypeName,
		searchattribute.StartTime:         request.StartTime,
		searchattribute.ExecutionTime:     request.ExecutionTime,
		searchattribute.ExecutionStatus:   request.Status.String(),
		searchattribute.TaskQueue:         request.TaskQueue,
	}

	if len(request.Memo.GetData()) > 0 {
		doc[searchattribute.Memo] = request.Memo.GetData()
		doc[searchattribute.MemoEncoding] = request.Memo.GetEncodingType().String()
	}

	typeMap, err := s.searchAttributesProvider.GetSearchAttributes(s.index, false)
	if err != nil {
		s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentGenerateFailuresCount)
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to read search attribute types: %v", err))
	}

	searchAttributes, err := searchattribute.Decode(request.SearchAttributes, &typeMap)
	if err != nil {
		s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentGenerateFailuresCount)
		return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to decode search attributes: %v", err))
	}
	for saName, saValue := range searchAttributes {
		doc[saName] = saValue
	}

	return doc, nil
}

func (s *visibilityStore) parseESDoc(hit *elastic.SearchHit, saTypeMap searchattribute.NameTypeMap, namespace string) (*store.InternalWorkflowExecutionInfo, error) {
	logParseError := func(fieldName string, fieldValue interface{}, err error, docID string) error {
		s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentParseFailuresCount)
		return serviceerror.NewUnavailable(fmt.Sprintf("Unable to parse Elasticsearch document(%s) %q field value %q: %v", docID, fieldName, fieldValue, err))
	}

	var sourceMap map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(hit.Source))
	// Very important line. See finishParseJSONValue bellow.
	d.UseNumber()
	if err := d.Decode(&sourceMap); err != nil {
		s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentParseFailuresCount)
		return nil, serviceerror.NewInternal(fmt.Sprintf("Unable to unmarshal JSON from Elasticsearch document(%s): %v", hit.Id, err))
	}

	var (
		isValidType            bool
		memo                   []byte
		memoEncoding           string
		customSearchAttributes map[string]interface{}
	)
	record := &store.InternalWorkflowExecutionInfo{}
	for fieldName, fieldValue := range sourceMap {
		switch fieldName {
		case searchattribute.NamespaceID,
			searchattribute.ExecutionDuration,
			searchattribute.VisibilityTaskKey:
			// Ignore these fields.
			continue
		case searchattribute.Memo:
			var memoStr string
			if memoStr, isValidType = fieldValue.(string); !isValidType {
				return nil, logParseError(fieldName, fieldValue, fmt.Errorf("%w: expected string got %T", errUnexpectedJSONFieldType, fieldValue), hit.Id)
			}
			var err error
			if memo, err = base64.StdEncoding.DecodeString(memoStr); err != nil {
				return nil, logParseError(fieldName, memoStr[:10], err, hit.Id)
			}
			continue
		case searchattribute.MemoEncoding:
			if memoEncoding, isValidType = fieldValue.(string); !isValidType {
				return nil, logParseError(fieldName, fieldValue, fmt.Errorf("%w: expected string got %T", errUnexpectedJSONFieldType, fieldValue), hit.Id)
			}
			continue
		}

		fieldType, err := saTypeMap.GetType(fieldName)
		if err != nil {
			// Silently ignore ErrInvalidName because it indicates unknown field in Elasticsearch document.
			if errors.Is(err, searchattribute.ErrInvalidName) {
				continue
			}
			s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentParseFailuresCount)
			return nil, serviceerror.NewUnavailable(fmt.Sprintf("Unable to get type for Elasticsearch document(%s) field %q: %v", hit.Id, fieldName, err))
		}

		fieldValueParsed, err := finishParseJSONValue(fieldValue, fieldType)
		if err != nil {
			return nil, logParseError(fieldName, fieldValue, err, hit.Id)
		}

		switch fieldName {
		case searchattribute.WorkflowID:
			record.WorkflowID = fieldValueParsed.(string)
		case searchattribute.RunID:
			record.RunID = fieldValueParsed.(string)
		case searchattribute.WorkflowType:
			record.TypeName = fieldValue.(string)
		case searchattribute.StartTime:
			record.StartTime = fieldValueParsed.(time.Time)
		case searchattribute.ExecutionTime:
			record.ExecutionTime = fieldValueParsed.(time.Time)
		case searchattribute.CloseTime:
			record.CloseTime = fieldValueParsed.(time.Time)
		case searchattribute.TaskQueue:
			record.TaskQueue = fieldValueParsed.(string)
		case searchattribute.ExecutionStatus:
			record.Status = enumspb.WorkflowExecutionStatus(enumspb.WorkflowExecutionStatus_value[fieldValueParsed.(string)])
		case searchattribute.HistoryLength:
			record.HistoryLength = fieldValueParsed.(int64)
		case searchattribute.StateTransitionCount:
			record.StateTransitionCount = fieldValueParsed.(int64)
		default:
			// All custom and predefined search attributes are handled here.
			if customSearchAttributes == nil {
				customSearchAttributes = map[string]interface{}{}
			}
			customSearchAttributes[fieldName] = fieldValueParsed
		}
	}

	if customSearchAttributes != nil {
		var err error
		record.SearchAttributes, err = searchattribute.Encode(customSearchAttributes, &saTypeMap)
		if err != nil {
			s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentParseFailuresCount)
			return nil, serviceerror.NewInternal(fmt.Sprintf("Unable to encode custom search attributes of Elasticsearch document(%s): %v", hit.Id, err))
		}
		err = searchattribute.ApplyAliases(s.searchAttributesMapper, record.SearchAttributes, namespace)
		if err != nil {
			return nil, err
		}
	}

	if memoEncoding != "" {
		record.Memo = persistence.NewDataBlob(memo, memoEncoding)
	} else if memo != nil {
		s.metricsClient.IncCounter(metrics.ElasticsearchVisibility, metrics.ElasticsearchDocumentParseFailuresCount)
		return nil, serviceerror.NewInternal(fmt.Sprintf("%q field is missing in Elasticsearch document(%s)", searchattribute.MemoEncoding, hit.Id))
	}

	return record, nil
}

// finishParseJSONValue finishes JSON parsing after json.Decode.
// json.Decode returns:
//     bool, for JSON booleans
//     json.Number, for JSON numbers (because of d.UseNumber())
//     string, for JSON strings
//     []interface{}, for JSON arrays
//     map[string]interface{}, for JSON objects (should never be a case)
//     nil for JSON null
func finishParseJSONValue(val interface{}, t enumspb.IndexedValueType) (interface{}, error) {
	// Custom search attributes support array of particular type.
	if arrayValue, isArray := val.([]interface{}); isArray {
		retArray := make([]interface{}, len(arrayValue))
		var lastErr error
		for i := 0; i < len(retArray); i++ {
			retArray[i], lastErr = finishParseJSONValue(arrayValue[i], t)
		}
		return retArray, lastErr
	}

	switch t {
	case enumspb.INDEXED_VALUE_TYPE_STRING, enumspb.INDEXED_VALUE_TYPE_KEYWORD, enumspb.INDEXED_VALUE_TYPE_DATETIME:
		stringVal, isString := val.(string)
		if !isString {
			return nil, fmt.Errorf("%w: expected string got %T", errUnexpectedJSONFieldType, val)
		}
		if t == enumspb.INDEXED_VALUE_TYPE_DATETIME {
			return time.Parse(time.RFC3339Nano, stringVal)
		}
		return stringVal, nil
	case enumspb.INDEXED_VALUE_TYPE_INT, enumspb.INDEXED_VALUE_TYPE_DOUBLE:
		numberVal, isNumber := val.(json.Number)
		if !isNumber {
			return nil, fmt.Errorf("%w: expected json.Number got %T", errUnexpectedJSONFieldType, val)
		}
		if t == enumspb.INDEXED_VALUE_TYPE_INT {
			return numberVal.Int64()
		}
		return numberVal.Float64()
	case enumspb.INDEXED_VALUE_TYPE_BOOL:
		boolVal, isBool := val.(bool)
		if !isBool {
			return nil, fmt.Errorf("%w: expected bool got %T", errUnexpectedJSONFieldType, val)
		}
		return boolVal, nil
	}

	panic(fmt.Sprintf("Unknown field type: %v", t))
}

func detailedErrorMessage(err error) string {
	var elasticErr *elastic.Error
	if !errors.As(err, &elasticErr) ||
		elasticErr.Details == nil ||
		len(elasticErr.Details.RootCause) == 0 ||
		(len(elasticErr.Details.RootCause) == 1 && elasticErr.Details.RootCause[0].Reason == elasticErr.Details.Reason) {
		return err.Error()
	}

	var sb strings.Builder
	sb.WriteString(elasticErr.Error())
	sb.WriteString(", root causes:")
	for i, rootCause := range elasticErr.Details.RootCause {
		sb.WriteString(fmt.Sprintf(" %s [type=%s]", rootCause.Reason, rootCause.Type))
		if i != len(elasticErr.Details.RootCause)-1 {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}
