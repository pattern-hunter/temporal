// The MIT License
//
// Copyright (c) 2021 Temporal Technologies, Inc.
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
// source: temporal/server/api/enums/v1/predicate.proto

package enums

import (
	reflect "reflect"
	"strconv"
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

type PredicateType int32

const (
	PREDICATE_TYPE_UNSPECIFIED  PredicateType = 0
	PREDICATE_TYPE_UNIVERSAL    PredicateType = 1
	PREDICATE_TYPE_EMPTY        PredicateType = 2
	PREDICATE_TYPE_AND          PredicateType = 3
	PREDICATE_TYPE_OR           PredicateType = 4
	PREDICATE_TYPE_NOT          PredicateType = 5
	PREDICATE_TYPE_NAMESPACE_ID PredicateType = 6
	PREDICATE_TYPE_TASK_TYPE    PredicateType = 7
)

// Enum value maps for PredicateType.
var (
	PredicateType_name = map[int32]string{
		0: "PREDICATE_TYPE_UNSPECIFIED",
		1: "PREDICATE_TYPE_UNIVERSAL",
		2: "PREDICATE_TYPE_EMPTY",
		3: "PREDICATE_TYPE_AND",
		4: "PREDICATE_TYPE_OR",
		5: "PREDICATE_TYPE_NOT",
		6: "PREDICATE_TYPE_NAMESPACE_ID",
		7: "PREDICATE_TYPE_TASK_TYPE",
	}
	PredicateType_value = map[string]int32{
		"PREDICATE_TYPE_UNSPECIFIED":  0,
		"PREDICATE_TYPE_UNIVERSAL":    1,
		"PREDICATE_TYPE_EMPTY":        2,
		"PREDICATE_TYPE_AND":          3,
		"PREDICATE_TYPE_OR":           4,
		"PREDICATE_TYPE_NOT":          5,
		"PREDICATE_TYPE_NAMESPACE_ID": 6,
		"PREDICATE_TYPE_TASK_TYPE":    7,
	}
)

func (x PredicateType) Enum() *PredicateType {
	p := new(PredicateType)
	*p = x
	return p
}

func (x PredicateType) String() string {
	switch x {
	case PREDICATE_TYPE_UNSPECIFIED:
		return "Unspecified"
	case PREDICATE_TYPE_UNIVERSAL:
		return "Universal"
	case PREDICATE_TYPE_EMPTY:
		return "Empty"
	case PREDICATE_TYPE_AND:
		return "And"
	case PREDICATE_TYPE_OR:
		return "Or"
	case PREDICATE_TYPE_NOT:
		return "Not"
	case PREDICATE_TYPE_NAMESPACE_ID:
		return "NamespaceId"
	case PREDICATE_TYPE_TASK_TYPE:
		return "TaskType"
	default:
		return strconv.Itoa(int(x))
	}

}

func (PredicateType) Descriptor() protoreflect.EnumDescriptor {
	return file_temporal_server_api_enums_v1_predicate_proto_enumTypes[0].Descriptor()
}

func (PredicateType) Type() protoreflect.EnumType {
	return &file_temporal_server_api_enums_v1_predicate_proto_enumTypes[0]
}

func (x PredicateType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PredicateType.Descriptor instead.
func (PredicateType) EnumDescriptor() ([]byte, []int) {
	return file_temporal_server_api_enums_v1_predicate_proto_rawDescGZIP(), []int{0}
}

var File_temporal_server_api_enums_v1_predicate_proto protoreflect.FileDescriptor

var file_temporal_server_api_enums_v1_predicate_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c,
	0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x76, 0x31, 0x2a, 0xed, 0x01, 0x0a,
	0x0d, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e,
	0x0a, 0x1a, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c,
	0x0a, 0x18, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x49, 0x56, 0x45, 0x52, 0x53, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14,
	0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45,
	0x4d, 0x50, 0x54, 0x59, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43,
	0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x4e, 0x44, 0x10, 0x03, 0x12, 0x15,
	0x0a, 0x11, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x4f, 0x52, 0x10, 0x04, 0x12, 0x16, 0x0a, 0x12, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41,
	0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x10, 0x05, 0x12, 0x1f, 0x0a,
	0x1b, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x49, 0x44, 0x10, 0x06, 0x12, 0x1c,
	0x0a, 0x18, 0x50, 0x52, 0x45, 0x44, 0x49, 0x43, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x07, 0x42, 0x2a, 0x5a, 0x28,
	0x67, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2f,
	0x76, 0x31, 0x3b, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_temporal_server_api_enums_v1_predicate_proto_rawDescOnce sync.Once
	file_temporal_server_api_enums_v1_predicate_proto_rawDescData = file_temporal_server_api_enums_v1_predicate_proto_rawDesc
)

func file_temporal_server_api_enums_v1_predicate_proto_rawDescGZIP() []byte {
	file_temporal_server_api_enums_v1_predicate_proto_rawDescOnce.Do(func() {
		file_temporal_server_api_enums_v1_predicate_proto_rawDescData = protoimpl.X.CompressGZIP(file_temporal_server_api_enums_v1_predicate_proto_rawDescData)
	})
	return file_temporal_server_api_enums_v1_predicate_proto_rawDescData
}

var file_temporal_server_api_enums_v1_predicate_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_temporal_server_api_enums_v1_predicate_proto_goTypes = []interface{}{
	(PredicateType)(0), // 0: temporal.server.api.enums.v1.PredicateType
}
var file_temporal_server_api_enums_v1_predicate_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_server_api_enums_v1_predicate_proto_init() }
func file_temporal_server_api_enums_v1_predicate_proto_init() {
	if File_temporal_server_api_enums_v1_predicate_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_temporal_server_api_enums_v1_predicate_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_server_api_enums_v1_predicate_proto_goTypes,
		DependencyIndexes: file_temporal_server_api_enums_v1_predicate_proto_depIdxs,
		EnumInfos:         file_temporal_server_api_enums_v1_predicate_proto_enumTypes,
	}.Build()
	File_temporal_server_api_enums_v1_predicate_proto = out.File
	file_temporal_server_api_enums_v1_predicate_proto_rawDesc = nil
	file_temporal_server_api_enums_v1_predicate_proto_goTypes = nil
	file_temporal_server_api_enums_v1_predicate_proto_depIdxs = nil
}
