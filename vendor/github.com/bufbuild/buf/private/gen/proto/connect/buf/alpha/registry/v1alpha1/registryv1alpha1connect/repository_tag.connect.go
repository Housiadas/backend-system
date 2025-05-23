// Copyright 2020-2025 Buf Technologies, Inc.
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

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: buf/alpha/registry/v1alpha1/repository_tag.proto

package registryv1alpha1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// RepositoryTagServiceName is the fully-qualified name of the RepositoryTagService service.
	RepositoryTagServiceName = "buf.alpha.registry.v1alpha1.RepositoryTagService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// RepositoryTagServiceGetRepositoryTagProcedure is the fully-qualified name of the
	// RepositoryTagService's GetRepositoryTag RPC.
	RepositoryTagServiceGetRepositoryTagProcedure = "/buf.alpha.registry.v1alpha1.RepositoryTagService/GetRepositoryTag"
	// RepositoryTagServiceCreateRepositoryTagProcedure is the fully-qualified name of the
	// RepositoryTagService's CreateRepositoryTag RPC.
	RepositoryTagServiceCreateRepositoryTagProcedure = "/buf.alpha.registry.v1alpha1.RepositoryTagService/CreateRepositoryTag"
	// RepositoryTagServiceUpdateRepositoryTagProcedure is the fully-qualified name of the
	// RepositoryTagService's UpdateRepositoryTag RPC.
	RepositoryTagServiceUpdateRepositoryTagProcedure = "/buf.alpha.registry.v1alpha1.RepositoryTagService/UpdateRepositoryTag"
	// RepositoryTagServiceListRepositoryTagsProcedure is the fully-qualified name of the
	// RepositoryTagService's ListRepositoryTags RPC.
	RepositoryTagServiceListRepositoryTagsProcedure = "/buf.alpha.registry.v1alpha1.RepositoryTagService/ListRepositoryTags"
	// RepositoryTagServiceListRepositoryTagsForReferenceProcedure is the fully-qualified name of the
	// RepositoryTagService's ListRepositoryTagsForReference RPC.
	RepositoryTagServiceListRepositoryTagsForReferenceProcedure = "/buf.alpha.registry.v1alpha1.RepositoryTagService/ListRepositoryTagsForReference"
)

// RepositoryTagServiceClient is a client for the buf.alpha.registry.v1alpha1.RepositoryTagService
// service.
type RepositoryTagServiceClient interface {
	// GetRepositoryTag gets an existing repository tag.
	GetRepositoryTag(context.Context, *connect.Request[v1alpha1.GetRepositoryTagRequest]) (*connect.Response[v1alpha1.GetRepositoryTagResponse], error)
	// CreateRepositoryTag creates a new repository tag.
	CreateRepositoryTag(context.Context, *connect.Request[v1alpha1.CreateRepositoryTagRequest]) (*connect.Response[v1alpha1.CreateRepositoryTagResponse], error)
	// UpdateRepositoryTag update an existing repository tag.
	UpdateRepositoryTag(context.Context, *connect.Request[v1alpha1.UpdateRepositoryTagRequest]) (*connect.Response[v1alpha1.UpdateRepositoryTagResponse], error)
	// ListRepositoryTags lists the repository tags associated with a Repository.
	ListRepositoryTags(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsResponse], error)
	// ListRepositoryTagsForReference lists the repository tags associated with a repository
	// reference name.
	ListRepositoryTagsForReference(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsForReferenceResponse], error)
}

// NewRepositoryTagServiceClient constructs a client for the
// buf.alpha.registry.v1alpha1.RepositoryTagService service. By default, it uses the Connect
// protocol with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed
// requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewRepositoryTagServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) RepositoryTagServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	repositoryTagServiceMethods := v1alpha1.File_buf_alpha_registry_v1alpha1_repository_tag_proto.Services().ByName("RepositoryTagService").Methods()
	return &repositoryTagServiceClient{
		getRepositoryTag: connect.NewClient[v1alpha1.GetRepositoryTagRequest, v1alpha1.GetRepositoryTagResponse](
			httpClient,
			baseURL+RepositoryTagServiceGetRepositoryTagProcedure,
			connect.WithSchema(repositoryTagServiceMethods.ByName("GetRepositoryTag")),
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
		createRepositoryTag: connect.NewClient[v1alpha1.CreateRepositoryTagRequest, v1alpha1.CreateRepositoryTagResponse](
			httpClient,
			baseURL+RepositoryTagServiceCreateRepositoryTagProcedure,
			connect.WithSchema(repositoryTagServiceMethods.ByName("CreateRepositoryTag")),
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
		updateRepositoryTag: connect.NewClient[v1alpha1.UpdateRepositoryTagRequest, v1alpha1.UpdateRepositoryTagResponse](
			httpClient,
			baseURL+RepositoryTagServiceUpdateRepositoryTagProcedure,
			connect.WithSchema(repositoryTagServiceMethods.ByName("UpdateRepositoryTag")),
			connect.WithIdempotency(connect.IdempotencyIdempotent),
			connect.WithClientOptions(opts...),
		),
		listRepositoryTags: connect.NewClient[v1alpha1.ListRepositoryTagsRequest, v1alpha1.ListRepositoryTagsResponse](
			httpClient,
			baseURL+RepositoryTagServiceListRepositoryTagsProcedure,
			connect.WithSchema(repositoryTagServiceMethods.ByName("ListRepositoryTags")),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
		listRepositoryTagsForReference: connect.NewClient[v1alpha1.ListRepositoryTagsForReferenceRequest, v1alpha1.ListRepositoryTagsForReferenceResponse](
			httpClient,
			baseURL+RepositoryTagServiceListRepositoryTagsForReferenceProcedure,
			connect.WithSchema(repositoryTagServiceMethods.ByName("ListRepositoryTagsForReference")),
			connect.WithIdempotency(connect.IdempotencyNoSideEffects),
			connect.WithClientOptions(opts...),
		),
	}
}

// repositoryTagServiceClient implements RepositoryTagServiceClient.
type repositoryTagServiceClient struct {
	getRepositoryTag               *connect.Client[v1alpha1.GetRepositoryTagRequest, v1alpha1.GetRepositoryTagResponse]
	createRepositoryTag            *connect.Client[v1alpha1.CreateRepositoryTagRequest, v1alpha1.CreateRepositoryTagResponse]
	updateRepositoryTag            *connect.Client[v1alpha1.UpdateRepositoryTagRequest, v1alpha1.UpdateRepositoryTagResponse]
	listRepositoryTags             *connect.Client[v1alpha1.ListRepositoryTagsRequest, v1alpha1.ListRepositoryTagsResponse]
	listRepositoryTagsForReference *connect.Client[v1alpha1.ListRepositoryTagsForReferenceRequest, v1alpha1.ListRepositoryTagsForReferenceResponse]
}

// GetRepositoryTag calls buf.alpha.registry.v1alpha1.RepositoryTagService.GetRepositoryTag.
func (c *repositoryTagServiceClient) GetRepositoryTag(ctx context.Context, req *connect.Request[v1alpha1.GetRepositoryTagRequest]) (*connect.Response[v1alpha1.GetRepositoryTagResponse], error) {
	return c.getRepositoryTag.CallUnary(ctx, req)
}

// CreateRepositoryTag calls buf.alpha.registry.v1alpha1.RepositoryTagService.CreateRepositoryTag.
func (c *repositoryTagServiceClient) CreateRepositoryTag(ctx context.Context, req *connect.Request[v1alpha1.CreateRepositoryTagRequest]) (*connect.Response[v1alpha1.CreateRepositoryTagResponse], error) {
	return c.createRepositoryTag.CallUnary(ctx, req)
}

// UpdateRepositoryTag calls buf.alpha.registry.v1alpha1.RepositoryTagService.UpdateRepositoryTag.
func (c *repositoryTagServiceClient) UpdateRepositoryTag(ctx context.Context, req *connect.Request[v1alpha1.UpdateRepositoryTagRequest]) (*connect.Response[v1alpha1.UpdateRepositoryTagResponse], error) {
	return c.updateRepositoryTag.CallUnary(ctx, req)
}

// ListRepositoryTags calls buf.alpha.registry.v1alpha1.RepositoryTagService.ListRepositoryTags.
func (c *repositoryTagServiceClient) ListRepositoryTags(ctx context.Context, req *connect.Request[v1alpha1.ListRepositoryTagsRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsResponse], error) {
	return c.listRepositoryTags.CallUnary(ctx, req)
}

// ListRepositoryTagsForReference calls
// buf.alpha.registry.v1alpha1.RepositoryTagService.ListRepositoryTagsForReference.
func (c *repositoryTagServiceClient) ListRepositoryTagsForReference(ctx context.Context, req *connect.Request[v1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsForReferenceResponse], error) {
	return c.listRepositoryTagsForReference.CallUnary(ctx, req)
}

// RepositoryTagServiceHandler is an implementation of the
// buf.alpha.registry.v1alpha1.RepositoryTagService service.
type RepositoryTagServiceHandler interface {
	// GetRepositoryTag gets an existing repository tag.
	GetRepositoryTag(context.Context, *connect.Request[v1alpha1.GetRepositoryTagRequest]) (*connect.Response[v1alpha1.GetRepositoryTagResponse], error)
	// CreateRepositoryTag creates a new repository tag.
	CreateRepositoryTag(context.Context, *connect.Request[v1alpha1.CreateRepositoryTagRequest]) (*connect.Response[v1alpha1.CreateRepositoryTagResponse], error)
	// UpdateRepositoryTag update an existing repository tag.
	UpdateRepositoryTag(context.Context, *connect.Request[v1alpha1.UpdateRepositoryTagRequest]) (*connect.Response[v1alpha1.UpdateRepositoryTagResponse], error)
	// ListRepositoryTags lists the repository tags associated with a Repository.
	ListRepositoryTags(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsResponse], error)
	// ListRepositoryTagsForReference lists the repository tags associated with a repository
	// reference name.
	ListRepositoryTagsForReference(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsForReferenceResponse], error)
}

// NewRepositoryTagServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewRepositoryTagServiceHandler(svc RepositoryTagServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	repositoryTagServiceMethods := v1alpha1.File_buf_alpha_registry_v1alpha1_repository_tag_proto.Services().ByName("RepositoryTagService").Methods()
	repositoryTagServiceGetRepositoryTagHandler := connect.NewUnaryHandler(
		RepositoryTagServiceGetRepositoryTagProcedure,
		svc.GetRepositoryTag,
		connect.WithSchema(repositoryTagServiceMethods.ByName("GetRepositoryTag")),
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	repositoryTagServiceCreateRepositoryTagHandler := connect.NewUnaryHandler(
		RepositoryTagServiceCreateRepositoryTagProcedure,
		svc.CreateRepositoryTag,
		connect.WithSchema(repositoryTagServiceMethods.ByName("CreateRepositoryTag")),
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	repositoryTagServiceUpdateRepositoryTagHandler := connect.NewUnaryHandler(
		RepositoryTagServiceUpdateRepositoryTagProcedure,
		svc.UpdateRepositoryTag,
		connect.WithSchema(repositoryTagServiceMethods.ByName("UpdateRepositoryTag")),
		connect.WithIdempotency(connect.IdempotencyIdempotent),
		connect.WithHandlerOptions(opts...),
	)
	repositoryTagServiceListRepositoryTagsHandler := connect.NewUnaryHandler(
		RepositoryTagServiceListRepositoryTagsProcedure,
		svc.ListRepositoryTags,
		connect.WithSchema(repositoryTagServiceMethods.ByName("ListRepositoryTags")),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	repositoryTagServiceListRepositoryTagsForReferenceHandler := connect.NewUnaryHandler(
		RepositoryTagServiceListRepositoryTagsForReferenceProcedure,
		svc.ListRepositoryTagsForReference,
		connect.WithSchema(repositoryTagServiceMethods.ByName("ListRepositoryTagsForReference")),
		connect.WithIdempotency(connect.IdempotencyNoSideEffects),
		connect.WithHandlerOptions(opts...),
	)
	return "/buf.alpha.registry.v1alpha1.RepositoryTagService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case RepositoryTagServiceGetRepositoryTagProcedure:
			repositoryTagServiceGetRepositoryTagHandler.ServeHTTP(w, r)
		case RepositoryTagServiceCreateRepositoryTagProcedure:
			repositoryTagServiceCreateRepositoryTagHandler.ServeHTTP(w, r)
		case RepositoryTagServiceUpdateRepositoryTagProcedure:
			repositoryTagServiceUpdateRepositoryTagHandler.ServeHTTP(w, r)
		case RepositoryTagServiceListRepositoryTagsProcedure:
			repositoryTagServiceListRepositoryTagsHandler.ServeHTTP(w, r)
		case RepositoryTagServiceListRepositoryTagsForReferenceProcedure:
			repositoryTagServiceListRepositoryTagsForReferenceHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedRepositoryTagServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedRepositoryTagServiceHandler struct{}

func (UnimplementedRepositoryTagServiceHandler) GetRepositoryTag(context.Context, *connect.Request[v1alpha1.GetRepositoryTagRequest]) (*connect.Response[v1alpha1.GetRepositoryTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryTagService.GetRepositoryTag is not implemented"))
}

func (UnimplementedRepositoryTagServiceHandler) CreateRepositoryTag(context.Context, *connect.Request[v1alpha1.CreateRepositoryTagRequest]) (*connect.Response[v1alpha1.CreateRepositoryTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryTagService.CreateRepositoryTag is not implemented"))
}

func (UnimplementedRepositoryTagServiceHandler) UpdateRepositoryTag(context.Context, *connect.Request[v1alpha1.UpdateRepositoryTagRequest]) (*connect.Response[v1alpha1.UpdateRepositoryTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryTagService.UpdateRepositoryTag is not implemented"))
}

func (UnimplementedRepositoryTagServiceHandler) ListRepositoryTags(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryTagService.ListRepositoryTags is not implemented"))
}

func (UnimplementedRepositoryTagServiceHandler) ListRepositoryTagsForReference(context.Context, *connect.Request[v1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[v1alpha1.ListRepositoryTagsForReferenceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("buf.alpha.registry.v1alpha1.RepositoryTagService.ListRepositoryTagsForReference is not implemented"))
}
