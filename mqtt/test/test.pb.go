//
// mqtt handler - test/test.proto
// Copyright (c) 2019 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
// Author: Matthias Schiffer and the Energy Manager development team
//
// This software code contained herein is licensed under the terms and conditions of
// the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
// You will find the corresponding license text in the LICENSE file.
// In case of any license issues please contact license@tq-group.com.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v4.25.3
// source: test/test.proto

package test

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Test struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	MessageCounter uint64                 `protobuf:"varint,1,opt,name=messageCounter,proto3" json:"messageCounter,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Test) Reset() {
	*x = Test{}
	mi := &file_test_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Test) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Test) ProtoMessage() {}

func (x *Test) ProtoReflect() protoreflect.Message {
	mi := &file_test_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Test.ProtoReflect.Descriptor instead.
func (*Test) Descriptor() ([]byte, []int) {
	return file_test_test_proto_rawDescGZIP(), []int{0}
}

func (x *Test) GetMessageCounter() uint64 {
	if x != nil {
		return x.MessageCounter
	}
	return 0
}

var File_test_test_proto protoreflect.FileDescriptor

var file_test_test_proto_rawDesc = string([]byte{
	0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2e, 0x0a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_test_test_proto_rawDescOnce sync.Once
	file_test_test_proto_rawDescData []byte
)

func file_test_test_proto_rawDescGZIP() []byte {
	file_test_test_proto_rawDescOnce.Do(func() {
		file_test_test_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_test_test_proto_rawDesc), len(file_test_test_proto_rawDesc)))
	})
	return file_test_test_proto_rawDescData
}

var file_test_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_test_proto_goTypes = []any{
	(*Test)(nil), // 0: Test
}
var file_test_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_test_proto_init() }
func file_test_test_proto_init() {
	if File_test_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_test_test_proto_rawDesc), len(file_test_test_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_test_proto_goTypes,
		DependencyIndexes: file_test_test_proto_depIdxs,
		MessageInfos:      file_test_test_proto_msgTypes,
	}.Build()
	File_test_test_proto = out.File
	file_test_test_proto_goTypes = nil
	file_test_test_proto_depIdxs = nil
}
