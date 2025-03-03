// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package channel

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ChannelAPIClient is the client API for ChannelAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelAPIClient interface {
	// Creates a new subscriber channel
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error)
	// Updates an existing channel resource
	UpdateChannel(ctx context.Context, in *UpdateChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Removes a subscribers channel
	DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Retrieves a collection of channels resource
	ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error)
	// Searches for channels
	SearchChannels(ctx context.Context, in *SearchChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error)
	// Retrieves a single channel resource
	GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*Channel, error)
	// Increment subscribers by one.
	IncrementSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Decrement subscribers by one.
	DecrementSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type channelAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelAPIClient(cc grpc.ClientConnInterface) ChannelAPIClient {
	return &channelAPIClient{cc}
}

func (c *channelAPIClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*CreateChannelResponse, error) {
	out := new(CreateChannelResponse)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/CreateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) UpdateChannel(ctx context.Context, in *UpdateChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/UpdateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/DeleteChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error) {
	out := new(ListChannelsResponse)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/ListChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) SearchChannels(ctx context.Context, in *SearchChannelsRequest, opts ...grpc.CallOption) (*ListChannelsResponse, error) {
	out := new(ListChannelsResponse)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/SearchChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) GetChannel(ctx context.Context, in *GetChannelRequest, opts ...grpc.CallOption) (*Channel, error) {
	out := new(Channel)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/GetChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) IncrementSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/IncrementSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelAPIClient) DecrementSubscribers(ctx context.Context, in *SubscribersRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gidyon.apis.ChannelAPI/DecrementSubscribers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelAPIServer is the server API for ChannelAPI service.
// All implementations must embed UnimplementedChannelAPIServer
// for forward compatibility
type ChannelAPIServer interface {
	// Creates a new subscriber channel
	CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error)
	// Updates an existing channel resource
	UpdateChannel(context.Context, *UpdateChannelRequest) (*empty.Empty, error)
	// Removes a subscribers channel
	DeleteChannel(context.Context, *DeleteChannelRequest) (*empty.Empty, error)
	// Retrieves a collection of channels resource
	ListChannels(context.Context, *ListChannelsRequest) (*ListChannelsResponse, error)
	// Searches for channels
	SearchChannels(context.Context, *SearchChannelsRequest) (*ListChannelsResponse, error)
	// Retrieves a single channel resource
	GetChannel(context.Context, *GetChannelRequest) (*Channel, error)
	// Increment subscribers by one.
	IncrementSubscribers(context.Context, *SubscribersRequest) (*empty.Empty, error)
	// Decrement subscribers by one.
	DecrementSubscribers(context.Context, *SubscribersRequest) (*empty.Empty, error)
	mustEmbedUnimplementedChannelAPIServer()
}

// UnimplementedChannelAPIServer must be embedded to have forward compatible implementations.
type UnimplementedChannelAPIServer struct {
}

func (UnimplementedChannelAPIServer) CreateChannel(context.Context, *CreateChannelRequest) (*CreateChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedChannelAPIServer) UpdateChannel(context.Context, *UpdateChannelRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChannel not implemented")
}
func (UnimplementedChannelAPIServer) DeleteChannel(context.Context, *DeleteChannelRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChannel not implemented")
}
func (UnimplementedChannelAPIServer) ListChannels(context.Context, *ListChannelsRequest) (*ListChannelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChannels not implemented")
}
func (UnimplementedChannelAPIServer) SearchChannels(context.Context, *SearchChannelsRequest) (*ListChannelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchChannels not implemented")
}
func (UnimplementedChannelAPIServer) GetChannel(context.Context, *GetChannelRequest) (*Channel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannel not implemented")
}
func (UnimplementedChannelAPIServer) IncrementSubscribers(context.Context, *SubscribersRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IncrementSubscribers not implemented")
}
func (UnimplementedChannelAPIServer) DecrementSubscribers(context.Context, *SubscribersRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecrementSubscribers not implemented")
}
func (UnimplementedChannelAPIServer) mustEmbedUnimplementedChannelAPIServer() {}

// UnsafeChannelAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelAPIServer will
// result in compilation errors.
type UnsafeChannelAPIServer interface {
	mustEmbedUnimplementedChannelAPIServer()
}

func RegisterChannelAPIServer(s grpc.ServiceRegistrar, srv ChannelAPIServer) {
	s.RegisterService(&_ChannelAPI_serviceDesc, srv)
}

func _ChannelAPI_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/CreateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).CreateChannel(ctx, req.(*CreateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_UpdateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).UpdateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/UpdateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).UpdateChannel(ctx, req.(*UpdateChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_DeleteChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).DeleteChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/DeleteChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).DeleteChannel(ctx, req.(*DeleteChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_ListChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).ListChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/ListChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).ListChannels(ctx, req.(*ListChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_SearchChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).SearchChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/SearchChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).SearchChannels(ctx, req.(*SearchChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_GetChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).GetChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/GetChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).GetChannel(ctx, req.(*GetChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_IncrementSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).IncrementSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/IncrementSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).IncrementSubscribers(ctx, req.(*SubscribersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelAPI_DecrementSubscribers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelAPIServer).DecrementSubscribers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gidyon.apis.ChannelAPI/DecrementSubscribers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelAPIServer).DecrementSubscribers(ctx, req.(*SubscribersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChannelAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gidyon.apis.ChannelAPI",
	HandlerType: (*ChannelAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChannel",
			Handler:    _ChannelAPI_CreateChannel_Handler,
		},
		{
			MethodName: "UpdateChannel",
			Handler:    _ChannelAPI_UpdateChannel_Handler,
		},
		{
			MethodName: "DeleteChannel",
			Handler:    _ChannelAPI_DeleteChannel_Handler,
		},
		{
			MethodName: "ListChannels",
			Handler:    _ChannelAPI_ListChannels_Handler,
		},
		{
			MethodName: "SearchChannels",
			Handler:    _ChannelAPI_SearchChannels_Handler,
		},
		{
			MethodName: "GetChannel",
			Handler:    _ChannelAPI_GetChannel_Handler,
		},
		{
			MethodName: "IncrementSubscribers",
			Handler:    _ChannelAPI_IncrementSubscribers_Handler,
		},
		{
			MethodName: "DecrementSubscribers",
			Handler:    _ChannelAPI_DecrementSubscribers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "channel.proto",
}
