// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: metricsapi.proto

package api

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MetricsApi_UpdateAllMetrics_FullMethodName = "/metricsapi.MetricsApi/UpdateAllMetrics"
	MetricsApi_GetPing_FullMethodName          = "/metricsapi.MetricsApi/GetPing"
)

// MetricsApiClient is the client API for MetricsApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsApiClient interface {
	UpdateAllMetrics(ctx context.Context, in *MetricsList, opts ...grpc.CallOption) (*UpdateResponse, error)
	GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Pong, error)
}

type metricsApiClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsApiClient(cc grpc.ClientConnInterface) MetricsApiClient {
	return &metricsApiClient{cc}
}

func (c *metricsApiClient) UpdateAllMetrics(ctx context.Context, in *MetricsList, opts ...grpc.CallOption) (*UpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, MetricsApi_UpdateAllMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsApiClient) GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Pong, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Pong)
	err := c.cc.Invoke(ctx, MetricsApi_GetPing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsApiServer is the server API for MetricsApi service.
// All implementations must embed UnimplementedMetricsApiServer
// for forward compatibility.
type MetricsApiServer interface {
	UpdateAllMetrics(context.Context, *MetricsList) (*UpdateResponse, error)
	GetPing(context.Context, *empty.Empty) (*Pong, error)
	mustEmbedUnimplementedMetricsApiServer()
}

// UnimplementedMetricsApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsApiServer struct{}

func (UnimplementedMetricsApiServer) UpdateAllMetrics(context.Context, *MetricsList) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAllMetrics not implemented")
}
func (UnimplementedMetricsApiServer) GetPing(context.Context, *empty.Empty) (*Pong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPing not implemented")
}
func (UnimplementedMetricsApiServer) mustEmbedUnimplementedMetricsApiServer() {}
func (UnimplementedMetricsApiServer) testEmbeddedByValue()                    {}

// UnsafeMetricsApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsApiServer will
// result in compilation errors.
type UnsafeMetricsApiServer interface {
	mustEmbedUnimplementedMetricsApiServer()
}

func RegisterMetricsApiServer(s grpc.ServiceRegistrar, srv MetricsApiServer) {
	// If the following call pancis, it indicates UnimplementedMetricsApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MetricsApi_ServiceDesc, srv)
}

func _MetricsApi_UpdateAllMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsApiServer).UpdateAllMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsApi_UpdateAllMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsApiServer).UpdateAllMetrics(ctx, req.(*MetricsList))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsApi_GetPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsApiServer).GetPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsApi_GetPing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsApiServer).GetPing(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsApi_ServiceDesc is the grpc.ServiceDesc for MetricsApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "metricsapi.MetricsApi",
	HandlerType: (*MetricsApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateAllMetrics",
			Handler:    _MetricsApi_UpdateAllMetrics_Handler,
		},
		{
			MethodName: "GetPing",
			Handler:    _MetricsApi_GetPing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metricsapi.proto",
}
