// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// plugins:
// 	protoc-gen-go
// 	protoc
// source: temporal/server/api/errordetails/v1/message.proto

// These error details extend failures defined in https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto

package errordetails

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TaskAlreadyStartedFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TaskAlreadyStartedFailure) Reset() {
	*x = TaskAlreadyStartedFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskAlreadyStartedFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskAlreadyStartedFailure) ProtoMessage() {}

func (x *TaskAlreadyStartedFailure) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskAlreadyStartedFailure.ProtoReflect.Descriptor instead.
func (*TaskAlreadyStartedFailure) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP(), []int{0}
}

type CurrentBranchChangedFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentBranchToken []byte `protobuf:"bytes,1,opt,name=current_branch_token,json=currentBranchToken,proto3" json:"current_branch_token,omitempty"`
	RequestBranchToken []byte `protobuf:"bytes,2,opt,name=request_branch_token,json=requestBranchToken,proto3" json:"request_branch_token,omitempty"`
}

func (x *CurrentBranchChangedFailure) Reset() {
	*x = CurrentBranchChangedFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CurrentBranchChangedFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CurrentBranchChangedFailure) ProtoMessage() {}

func (x *CurrentBranchChangedFailure) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CurrentBranchChangedFailure.ProtoReflect.Descriptor instead.
func (*CurrentBranchChangedFailure) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP(), []int{1}
}

func (x *CurrentBranchChangedFailure) GetCurrentBranchToken() []byte {
	if x != nil {
		return x.CurrentBranchToken
	}
	return nil
}

func (x *CurrentBranchChangedFailure) GetRequestBranchToken() []byte {
	if x != nil {
		return x.RequestBranchToken
	}
	return nil
}

type ShardOwnershipLostFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerHost   string `protobuf:"bytes,1,opt,name=owner_host,json=ownerHost,proto3" json:"owner_host,omitempty"`
	CurrentHost string `protobuf:"bytes,2,opt,name=current_host,json=currentHost,proto3" json:"current_host,omitempty"`
}

func (x *ShardOwnershipLostFailure) Reset() {
	*x = ShardOwnershipLostFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardOwnershipLostFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardOwnershipLostFailure) ProtoMessage() {}

func (x *ShardOwnershipLostFailure) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardOwnershipLostFailure.ProtoReflect.Descriptor instead.
func (*ShardOwnershipLostFailure) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP(), []int{2}
}

func (x *ShardOwnershipLostFailure) GetOwnerHost() string {
	if x != nil {
		return x.OwnerHost
	}
	return ""
}

func (x *ShardOwnershipLostFailure) GetCurrentHost() string {
	if x != nil {
		return x.CurrentHost
	}
	return ""
}

type RetryReplicationFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NamespaceId       string `protobuf:"bytes,1,opt,name=namespace_id,json=namespaceId,proto3" json:"namespace_id,omitempty"`
	WorkflowId        string `protobuf:"bytes,2,opt,name=workflow_id,json=workflowId,proto3" json:"workflow_id,omitempty"`
	RunId             string `protobuf:"bytes,3,opt,name=run_id,json=runId,proto3" json:"run_id,omitempty"`
	StartEventId      int64  `protobuf:"varint,4,opt,name=start_event_id,json=startEventId,proto3" json:"start_event_id,omitempty"`
	StartEventVersion int64  `protobuf:"varint,5,opt,name=start_event_version,json=startEventVersion,proto3" json:"start_event_version,omitempty"`
	EndEventId        int64  `protobuf:"varint,6,opt,name=end_event_id,json=endEventId,proto3" json:"end_event_id,omitempty"`
	EndEventVersion   int64  `protobuf:"varint,7,opt,name=end_event_version,json=endEventVersion,proto3" json:"end_event_version,omitempty"`
}

func (x *RetryReplicationFailure) Reset() {
	*x = RetryReplicationFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetryReplicationFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetryReplicationFailure) ProtoMessage() {}

func (x *RetryReplicationFailure) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetryReplicationFailure.ProtoReflect.Descriptor instead.
func (*RetryReplicationFailure) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP(), []int{3}
}

func (x *RetryReplicationFailure) GetNamespaceId() string {
	if x != nil {
		return x.NamespaceId
	}
	return ""
}

func (x *RetryReplicationFailure) GetWorkflowId() string {
	if x != nil {
		return x.WorkflowId
	}
	return ""
}

func (x *RetryReplicationFailure) GetRunId() string {
	if x != nil {
		return x.RunId
	}
	return ""
}

func (x *RetryReplicationFailure) GetStartEventId() int64 {
	if x != nil {
		return x.StartEventId
	}
	return 0
}

func (x *RetryReplicationFailure) GetStartEventVersion() int64 {
	if x != nil {
		return x.StartEventVersion
	}
	return 0
}

func (x *RetryReplicationFailure) GetEndEventId() int64 {
	if x != nil {
		return x.EndEventId
	}
	return 0
}

func (x *RetryReplicationFailure) GetEndEventVersion() int64 {
	if x != nil {
		return x.EndEventVersion
	}
	return 0
}

type StickyWorkerUnavailableFailure struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StickyWorkerUnavailableFailure) Reset() {
	*x = StickyWorkerUnavailableFailure{}
	if protoimpl.UnsafeEnabled {
		mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StickyWorkerUnavailableFailure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StickyWorkerUnavailableFailure) ProtoMessage() {}

func (x *StickyWorkerUnavailableFailure) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_errordetails_v1_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StickyWorkerUnavailableFailure.ProtoReflect.Descriptor instead.
func (*StickyWorkerUnavailableFailure) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP(), []int{4}
}

var File_temporal_server_api_errordetails_v1_message_proto protoreflect.FileDescriptor

var file_temporal_server_api_errordetails_v1_message_proto_rawDesc = []byte{
	0x0a, 0x31, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x23, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x1b, 0x0a, 0x19, 0x54, 0x61, 0x73, 0x6b,
	0x41, 0x6c, 0x72, 0x65, 0x61, 0x64, 0x79, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x22, 0x81, 0x01, 0x0a, 0x1b, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x42, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x12, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x30, 0x0a, 0x14, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x5d, 0x0a, 0x19, 0x53, 0x68, 0x61,
	0x72, 0x64, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x4c, 0x6f, 0x73, 0x74, 0x46,
	0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f,
	0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x5f, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x48, 0x6f, 0x73, 0x74, 0x22, 0x98, 0x02, 0x0a, 0x17, 0x52, 0x65, 0x74,
	0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69,
	0x6c, 0x75, 0x72, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x75, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x75, 0x6e, 0x49, 0x64, 0x12,
	0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0c, 0x65, 0x6e, 0x64, 0x5f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x6e, 0x64,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x65, 0x6e, 0x64, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0f, 0x65, 0x6e, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x20, 0x0a, 0x1e, 0x53, 0x74, 0x69, 0x63, 0x6b, 0x79, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x55, 0x6e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x46, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70,
	0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2f,
	0x76, 0x31, 0x3b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_temporal_server_api_errordetails_v1_message_proto_rawDescOnce sync.Once
	file_temporal_server_api_errordetails_v1_message_proto_rawDescData = file_temporal_server_api_errordetails_v1_message_proto_rawDesc
)

func file_temporal_server_api_errordetails_v1_message_proto_rawDescGZIP() []byte {
	file_temporal_server_api_errordetails_v1_message_proto_rawDescOnce.Do(func() {
		file_temporal_server_api_errordetails_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_temporal_server_api_errordetails_v1_message_proto_rawDescData)
	})
	return file_temporal_server_api_errordetails_v1_message_proto_rawDescData
}

var file_temporal_server_api_errordetails_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_temporal_server_api_errordetails_v1_message_proto_goTypes = []interface{}{
	(*TaskAlreadyStartedFailure)(nil),      // 0: temporal.server.api.errordetails.v1.TaskAlreadyStartedFailure
	(*CurrentBranchChangedFailure)(nil),    // 1: temporal.server.api.errordetails.v1.CurrentBranchChangedFailure
	(*ShardOwnershipLostFailure)(nil),      // 2: temporal.server.api.errordetails.v1.ShardOwnershipLostFailure
	(*RetryReplicationFailure)(nil),        // 3: temporal.server.api.errordetails.v1.RetryReplicationFailure
	(*StickyWorkerUnavailableFailure)(nil), // 4: temporal.server.api.errordetails.v1.StickyWorkerUnavailableFailure
}
var file_temporal_server_api_errordetails_v1_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_server_api_errordetails_v1_message_proto_init() }
func file_temporal_server_api_errordetails_v1_message_proto_init() {
	if File_temporal_server_api_errordetails_v1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_temporal_server_api_errordetails_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskAlreadyStartedFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_temporal_server_api_errordetails_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CurrentBranchChangedFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_temporal_server_api_errordetails_v1_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardOwnershipLostFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_temporal_server_api_errordetails_v1_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RetryReplicationFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_temporal_server_api_errordetails_v1_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StickyWorkerUnavailableFailure); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_temporal_server_api_errordetails_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_server_api_errordetails_v1_message_proto_goTypes,
		DependencyIndexes: file_temporal_server_api_errordetails_v1_message_proto_depIdxs,
		MessageInfos:      file_temporal_server_api_errordetails_v1_message_proto_msgTypes,
	}.Build()
	File_temporal_server_api_errordetails_v1_message_proto = out.File
	file_temporal_server_api_errordetails_v1_message_proto_rawDesc = nil
	file_temporal_server_api_errordetails_v1_message_proto_goTypes = nil
	file_temporal_server_api_errordetails_v1_message_proto_depIdxs = nil
}
