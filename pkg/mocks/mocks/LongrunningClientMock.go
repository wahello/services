// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	longrunning "github.com/gidyon/services/pkg/api/longrunning"

	mock "github.com/stretchr/testify/mock"
)

// LongrunningClientMock is an autogenerated mock type for the LongrunningClientMock type
type LongrunningClientMock struct {
	mock.Mock
}

// CreateOperation provides a mock function with given fields: ctx, in, opts
func (_m *LongrunningClientMock) CreateOperation(ctx context.Context, in *longrunning.CreateOperationRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *longrunning.Operation
	if rf, ok := ret.Get(0).(func(context.Context, *longrunning.CreateOperationRequest, ...grpc.CallOption) *longrunning.Operation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*longrunning.Operation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *longrunning.CreateOperationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOperation provides a mock function with given fields: ctx, in, opts
func (_m *LongrunningClientMock) DeleteOperation(ctx context.Context, in *longrunning.DeleteOperationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *longrunning.DeleteOperationRequest, ...grpc.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *longrunning.DeleteOperationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperation provides a mock function with given fields: ctx, in, opts
func (_m *LongrunningClientMock) GetOperation(ctx context.Context, in *longrunning.GetOperationRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *longrunning.Operation
	if rf, ok := ret.Get(0).(func(context.Context, *longrunning.GetOperationRequest, ...grpc.CallOption) *longrunning.Operation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*longrunning.Operation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *longrunning.GetOperationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOperations provides a mock function with given fields: ctx, in, opts
func (_m *LongrunningClientMock) ListOperations(ctx context.Context, in *longrunning.ListOperationsRequest, opts ...grpc.CallOption) (*longrunning.ListOperationsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *longrunning.ListOperationsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *longrunning.ListOperationsRequest, ...grpc.CallOption) *longrunning.ListOperationsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*longrunning.ListOperationsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *longrunning.ListOperationsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOperation provides a mock function with given fields: ctx, in, opts
func (_m *LongrunningClientMock) UpdateOperation(ctx context.Context, in *longrunning.UpdateOperationRequest, opts ...grpc.CallOption) (*longrunning.Operation, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *longrunning.Operation
	if rf, ok := ret.Get(0).(func(context.Context, *longrunning.UpdateOperationRequest, ...grpc.CallOption) *longrunning.Operation); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*longrunning.Operation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *longrunning.UpdateOperationRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
