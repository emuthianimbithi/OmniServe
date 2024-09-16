// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: omniserve_proto/omniserve.proto

package omniserve_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	OmniServe_PushFiles_FullMethodName = "/omniserve.OmniServe/PushFiles"
)

// OmniServeClient is the client API for OmniServe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OmniServeClient interface {
	PushFiles(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, PushResponse], error)
}

type omniServeClient struct {
	cc grpc.ClientConnInterface
}

func NewOmniServeClient(cc grpc.ClientConnInterface) OmniServeClient {
	return &omniServeClient{cc}
}

func (c *omniServeClient) PushFiles(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, PushResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &OmniServe_ServiceDesc.Streams[0], OmniServe_PushFiles_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileChunk, PushResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OmniServe_PushFilesClient = grpc.ClientStreamingClient[FileChunk, PushResponse]

// OmniServeServer is the server API for OmniServe service.
// All implementations must embed UnimplementedOmniServeServer
// for forward compatibility.
type OmniServeServer interface {
	PushFiles(grpc.ClientStreamingServer[FileChunk, PushResponse]) error
	mustEmbedUnimplementedOmniServeServer()
}

// UnimplementedOmniServeServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOmniServeServer struct{}

func (UnimplementedOmniServeServer) PushFiles(grpc.ClientStreamingServer[FileChunk, PushResponse]) error {
	return status.Errorf(codes.Unimplemented, "method PushFiles not implemented")
}
func (UnimplementedOmniServeServer) mustEmbedUnimplementedOmniServeServer() {}
func (UnimplementedOmniServeServer) testEmbeddedByValue()                   {}

// UnsafeOmniServeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OmniServeServer will
// result in compilation errors.
type UnsafeOmniServeServer interface {
	mustEmbedUnimplementedOmniServeServer()
}

func RegisterOmniServeServer(s grpc.ServiceRegistrar, srv OmniServeServer) {
	// If the following call pancis, it indicates UnimplementedOmniServeServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OmniServe_ServiceDesc, srv)
}

func _OmniServe_PushFiles_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OmniServeServer).PushFiles(&grpc.GenericServerStream[FileChunk, PushResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type OmniServe_PushFilesServer = grpc.ClientStreamingServer[FileChunk, PushResponse]

// OmniServe_ServiceDesc is the grpc.ServiceDesc for OmniServe service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OmniServe_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "omniserve.OmniServe",
	HandlerType: (*OmniServeServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushFiles",
			Handler:       _OmniServe_PushFiles_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "omniserve_proto/omniserve.proto",
}