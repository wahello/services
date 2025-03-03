// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	channel "github.com/gidyon/services/pkg/api/channel"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// ChannelAPIClientMock is an autogenerated mock type for the ChannelAPIClientMock type
type ChannelAPIClientMock struct {
	mock.Mock
}

// CreateChannel provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) CreateChannel(ctx context.Context, in *channel.CreateChannelRequest, opts ...grpc.CallOption) (*channel.CreateChannelResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *channel.CreateChannelResponse
	if rf, ok := ret.Get(0).(func(context.Context, *channel.CreateChannelRequest, ...grpc.CallOption) *channel.CreateChannelResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*channel.CreateChannelResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.CreateChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecrementSubscribers provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) DecrementSubscribers(ctx context.Context, in *channel.SubscribersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *channel.SubscribersRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.SubscribersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteChannel provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) DeleteChannel(ctx context.Context, in *channel.DeleteChannelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *channel.DeleteChannelRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.DeleteChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChannel provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) GetChannel(ctx context.Context, in *channel.GetChannelRequest, opts ...grpc.CallOption) (*channel.Channel, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *channel.Channel
	if rf, ok := ret.Get(0).(func(context.Context, *channel.GetChannelRequest, ...grpc.CallOption) *channel.Channel); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*channel.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.GetChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementSubscribers provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) IncrementSubscribers(ctx context.Context, in *channel.SubscribersRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *channel.SubscribersRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.SubscribersRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListChannels provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) ListChannels(ctx context.Context, in *channel.ListChannelsRequest, opts ...grpc.CallOption) (*channel.ListChannelsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *channel.ListChannelsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *channel.ListChannelsRequest, ...grpc.CallOption) *channel.ListChannelsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*channel.ListChannelsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.ListChannelsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchChannels provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) SearchChannels(ctx context.Context, in *channel.SearchChannelsRequest, opts ...grpc.CallOption) (*channel.ListChannelsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *channel.ListChannelsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *channel.SearchChannelsRequest, ...grpc.CallOption) *channel.ListChannelsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*channel.ListChannelsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.SearchChannelsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateChannel provides a mock function with given fields: ctx, in, opts
func (_m *ChannelAPIClientMock) UpdateChannel(ctx context.Context, in *channel.UpdateChannelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *channel.UpdateChannelRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *channel.UpdateChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
