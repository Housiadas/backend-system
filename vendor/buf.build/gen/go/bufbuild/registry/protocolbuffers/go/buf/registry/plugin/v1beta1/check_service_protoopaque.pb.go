// Copyright 2023-2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        (unknown)
// source: buf/registry/plugin/v1beta1/check_service.proto

//go:build protoopaque

package pluginv1beta1

import (
	v1 "buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go/buf/plugin/check/v1"
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListRulesRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_PageSize    uint32                 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3"`
	xxx_hidden_PageToken   string                 `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3"`
	xxx_hidden_ResourceRef *ResourceRef           `protobuf:"bytes,3,opt,name=resource_ref,json=resourceRef,proto3"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ListRulesRequest) Reset() {
	*x = ListRulesRequest{}
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRulesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRulesRequest) ProtoMessage() {}

func (x *ListRulesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ListRulesRequest) GetPageSize() uint32 {
	if x != nil {
		return x.xxx_hidden_PageSize
	}
	return 0
}

func (x *ListRulesRequest) GetPageToken() string {
	if x != nil {
		return x.xxx_hidden_PageToken
	}
	return ""
}

func (x *ListRulesRequest) GetResourceRef() *ResourceRef {
	if x != nil {
		return x.xxx_hidden_ResourceRef
	}
	return nil
}

func (x *ListRulesRequest) SetPageSize(v uint32) {
	x.xxx_hidden_PageSize = v
}

func (x *ListRulesRequest) SetPageToken(v string) {
	x.xxx_hidden_PageToken = v
}

func (x *ListRulesRequest) SetResourceRef(v *ResourceRef) {
	x.xxx_hidden_ResourceRef = v
}

func (x *ListRulesRequest) HasResourceRef() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_ResourceRef != nil
}

func (x *ListRulesRequest) ClearResourceRef() {
	x.xxx_hidden_ResourceRef = nil
}

type ListRulesRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// The maximum number of items to return.
	//
	// The default value is 10.
	PageSize uint32
	// The page to start from.
	//
	// If empty, the first page is returned.
	PageToken string
	// The reference to get Rules for.
	//
	// See the documentation on ResourceRef for resource resolution details.
	//
	// Once the resource is resolved, the following Rules are returned:
	//   - If a Plugin is referenced, the Rules of the the default Label are returned.
	//   - If a Label is referenced, the Rules of the Commit of this Label are returned.
	//   - If a Commit is referenced, the Rules for this Commit are returned.
	ResourceRef *ResourceRef
}

func (b0 ListRulesRequest_builder) Build() *ListRulesRequest {
	m0 := &ListRulesRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_PageSize = b.PageSize
	x.xxx_hidden_PageToken = b.PageToken
	x.xxx_hidden_ResourceRef = b.ResourceRef
	return m0
}

type ListRulesResponse struct {
	state                    protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_NextPageToken string                 `protobuf:"bytes,1,opt,name=next_page_token,json=nextPageToken,proto3"`
	xxx_hidden_Rules         *[]*v1.Rule            `protobuf:"bytes,2,rep,name=rules,proto3"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *ListRulesResponse) Reset() {
	*x = ListRulesResponse{}
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRulesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRulesResponse) ProtoMessage() {}

func (x *ListRulesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ListRulesResponse) GetNextPageToken() string {
	if x != nil {
		return x.xxx_hidden_NextPageToken
	}
	return ""
}

func (x *ListRulesResponse) GetRules() []*v1.Rule {
	if x != nil {
		if x.xxx_hidden_Rules != nil {
			return *x.xxx_hidden_Rules
		}
	}
	return nil
}

func (x *ListRulesResponse) SetNextPageToken(v string) {
	x.xxx_hidden_NextPageToken = v
}

func (x *ListRulesResponse) SetRules(v []*v1.Rule) {
	x.xxx_hidden_Rules = &v
}

type ListRulesResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// The next page token.
	//
	// If empty, there are no more pages.
	NextPageToken string
	// The listed Rules. This is part of the Bufplugin API, and can be returned from
	// a plugin's ListRules implementation.
	Rules []*v1.Rule
}

func (b0 ListRulesResponse_builder) Build() *ListRulesResponse {
	m0 := &ListRulesResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_NextPageToken = b.NextPageToken
	x.xxx_hidden_Rules = &b.Rules
	return m0
}

type ListCategoriesRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_PageSize    uint32                 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3"`
	xxx_hidden_PageToken   string                 `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3"`
	xxx_hidden_ResourceRef *ResourceRef           `protobuf:"bytes,3,opt,name=resource_ref,json=resourceRef,proto3"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ListCategoriesRequest) Reset() {
	*x = ListCategoriesRequest{}
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCategoriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCategoriesRequest) ProtoMessage() {}

func (x *ListCategoriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ListCategoriesRequest) GetPageSize() uint32 {
	if x != nil {
		return x.xxx_hidden_PageSize
	}
	return 0
}

func (x *ListCategoriesRequest) GetPageToken() string {
	if x != nil {
		return x.xxx_hidden_PageToken
	}
	return ""
}

func (x *ListCategoriesRequest) GetResourceRef() *ResourceRef {
	if x != nil {
		return x.xxx_hidden_ResourceRef
	}
	return nil
}

func (x *ListCategoriesRequest) SetPageSize(v uint32) {
	x.xxx_hidden_PageSize = v
}

func (x *ListCategoriesRequest) SetPageToken(v string) {
	x.xxx_hidden_PageToken = v
}

func (x *ListCategoriesRequest) SetResourceRef(v *ResourceRef) {
	x.xxx_hidden_ResourceRef = v
}

func (x *ListCategoriesRequest) HasResourceRef() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_ResourceRef != nil
}

func (x *ListCategoriesRequest) ClearResourceRef() {
	x.xxx_hidden_ResourceRef = nil
}

type ListCategoriesRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// The maximum number of items to return.
	//
	// The default value is 10.
	PageSize uint32
	// The page to start from.
	//
	// If empty, the first page is returned.
	PageToken string
	// The reference to get Categories for.
	//
	// See the documentation on ResourceRef for resource resolution details.
	//
	// Once the resource is resolved, the following Categories are returned:
	//   - If a Plugin is referenced, the Categories of the the default Label are returned.
	//   - If a Label is referenced, the Categories of the Commit of this Label are returned.
	//   - If a Commit is referenced, the Categories for this Commit are returned.
	ResourceRef *ResourceRef
}

func (b0 ListCategoriesRequest_builder) Build() *ListCategoriesRequest {
	m0 := &ListCategoriesRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_PageSize = b.PageSize
	x.xxx_hidden_PageToken = b.PageToken
	x.xxx_hidden_ResourceRef = b.ResourceRef
	return m0
}

type ListCategoriesResponse struct {
	state                    protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_NextPageToken string                 `protobuf:"bytes,1,opt,name=next_page_token,json=nextPageToken,proto3"`
	xxx_hidden_Categories    *[]*v1.Category        `protobuf:"bytes,2,rep,name=categories,proto3"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *ListCategoriesResponse) Reset() {
	*x = ListCategoriesResponse{}
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCategoriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCategoriesResponse) ProtoMessage() {}

func (x *ListCategoriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ListCategoriesResponse) GetNextPageToken() string {
	if x != nil {
		return x.xxx_hidden_NextPageToken
	}
	return ""
}

func (x *ListCategoriesResponse) GetCategories() []*v1.Category {
	if x != nil {
		if x.xxx_hidden_Categories != nil {
			return *x.xxx_hidden_Categories
		}
	}
	return nil
}

func (x *ListCategoriesResponse) SetNextPageToken(v string) {
	x.xxx_hidden_NextPageToken = v
}

func (x *ListCategoriesResponse) SetCategories(v []*v1.Category) {
	x.xxx_hidden_Categories = &v
}

type ListCategoriesResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// The next page token.
	//
	// If empty, there are no more pages.
	NextPageToken string
	// The listed Categories. This is part of the Bufplugin API, and can be returned from
	// a plugin's ListCategories implementation.
	Categories []*v1.Category
}

func (b0 ListCategoriesResponse_builder) Build() *ListCategoriesResponse {
	m0 := &ListCategoriesResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_NextPageToken = b.NextPageToken
	x.xxx_hidden_Categories = &b.Categories
	return m0
}

var File_buf_registry_plugin_v1beta1_check_service_proto protoreflect.FileDescriptor

var file_buf_registry_plugin_v1beta1_check_service_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1b, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x22,
	0x62, 0x75, 0x66, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2a, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb7, 0x01, 0x0a, 0x10,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x25, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x42, 0x08, 0xba, 0x48, 0x05, 0x2a, 0x03, 0x18, 0xfa, 0x01, 0x52, 0x08, 0x70,
	0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x27, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0x18, 0x80, 0x20, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x53, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x66,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x66,
	0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x52, 0x65, 0x66, 0x22, 0x76, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x0f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x18, 0x80, 0x20, 0x52, 0x0d, 0x6e,
	0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2f, 0x0a, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x62, 0x75,
	0x66, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x22, 0xbc, 0x01,
	0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x08, 0xba, 0x48, 0x05, 0x2a,
	0x03, 0x18, 0xfa, 0x01, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x27,
	0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x18, 0x80, 0x20, 0x52, 0x09, 0x70, 0x61,
	0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x53, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e,
	0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x66, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52,
	0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x66, 0x22, 0x89, 0x01, 0x0a,
	0x16, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f,
	0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0x18, 0x80, 0x20, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74,
	0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3d, 0x0a, 0x0a, 0x63, 0x61, 0x74,
	0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x62, 0x75, 0x66, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x32, 0xff, 0x01, 0x0a, 0x0c, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6f, 0x0a, 0x09, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x2d, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x03, 0x90, 0x02, 0x01, 0x12, 0x7e, 0x0a, 0x0e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x32, 0x2e, 0x62,
	0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x33, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x03, 0x90, 0x02, 0x01, 0x42, 0x61, 0x5a, 0x5f, 0x62, 0x75,
	0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x62,
	0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66, 0x66, 0x65, 0x72, 0x73,
	0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_buf_registry_plugin_v1beta1_check_service_proto_goTypes = []any{
	(*ListRulesRequest)(nil),       // 0: buf.registry.plugin.v1beta1.ListRulesRequest
	(*ListRulesResponse)(nil),      // 1: buf.registry.plugin.v1beta1.ListRulesResponse
	(*ListCategoriesRequest)(nil),  // 2: buf.registry.plugin.v1beta1.ListCategoriesRequest
	(*ListCategoriesResponse)(nil), // 3: buf.registry.plugin.v1beta1.ListCategoriesResponse
	(*ResourceRef)(nil),            // 4: buf.registry.plugin.v1beta1.ResourceRef
	(*v1.Rule)(nil),                // 5: buf.plugin.check.v1.Rule
	(*v1.Category)(nil),            // 6: buf.plugin.check.v1.Category
}
var file_buf_registry_plugin_v1beta1_check_service_proto_depIdxs = []int32{
	4, // 0: buf.registry.plugin.v1beta1.ListRulesRequest.resource_ref:type_name -> buf.registry.plugin.v1beta1.ResourceRef
	5, // 1: buf.registry.plugin.v1beta1.ListRulesResponse.rules:type_name -> buf.plugin.check.v1.Rule
	4, // 2: buf.registry.plugin.v1beta1.ListCategoriesRequest.resource_ref:type_name -> buf.registry.plugin.v1beta1.ResourceRef
	6, // 3: buf.registry.plugin.v1beta1.ListCategoriesResponse.categories:type_name -> buf.plugin.check.v1.Category
	0, // 4: buf.registry.plugin.v1beta1.CheckService.ListRules:input_type -> buf.registry.plugin.v1beta1.ListRulesRequest
	2, // 5: buf.registry.plugin.v1beta1.CheckService.ListCategories:input_type -> buf.registry.plugin.v1beta1.ListCategoriesRequest
	1, // 6: buf.registry.plugin.v1beta1.CheckService.ListRules:output_type -> buf.registry.plugin.v1beta1.ListRulesResponse
	3, // 7: buf.registry.plugin.v1beta1.CheckService.ListCategories:output_type -> buf.registry.plugin.v1beta1.ListCategoriesResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_buf_registry_plugin_v1beta1_check_service_proto_init() }
func file_buf_registry_plugin_v1beta1_check_service_proto_init() {
	if File_buf_registry_plugin_v1beta1_check_service_proto != nil {
		return
	}
	file_buf_registry_plugin_v1beta1_resource_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_buf_registry_plugin_v1beta1_check_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_buf_registry_plugin_v1beta1_check_service_proto_goTypes,
		DependencyIndexes: file_buf_registry_plugin_v1beta1_check_service_proto_depIdxs,
		MessageInfos:      file_buf_registry_plugin_v1beta1_check_service_proto_msgTypes,
	}.Build()
	File_buf_registry_plugin_v1beta1_check_service_proto = out.File
	file_buf_registry_plugin_v1beta1_check_service_proto_rawDesc = nil
	file_buf_registry_plugin_v1beta1_check_service_proto_goTypes = nil
	file_buf_registry_plugin_v1beta1_check_service_proto_depIdxs = nil
}
