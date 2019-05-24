// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import grpc "github.com/ixoja/shorten/vendor/google.golang.org/grpc"
import grpcapi "github.com/ixoja/shorten/internal/grpcapi"
import mock "github.com/stretchr/testify/mock"

// ShortenServiceClient is an autogenerated mock type for the ShortenServiceClient type
type ShortenServiceClient struct {
	mock.Mock
}

// RedirectURL provides a mock function with given fields: ctx, in, opts
func (_m *ShortenServiceClient) RedirectURL(ctx context.Context, in *grpcapi.RedirectURLRequest, opts ...grpc.CallOption) (*grpcapi.RedirectURLResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *grpcapi.RedirectURLResponse
	if rf, ok := ret.Get(0).(func(context.Context, *grpcapi.RedirectURLRequest, ...grpc.CallOption) *grpcapi.RedirectURLResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpcapi.RedirectURLResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpcapi.RedirectURLRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Shorten provides a mock function with given fields: ctx, in, opts
func (_m *ShortenServiceClient) Shorten(ctx context.Context, in *grpcapi.ShortenRequest, opts ...grpc.CallOption) (*grpcapi.ShortenResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *grpcapi.ShortenResponse
	if rf, ok := ret.Get(0).(func(context.Context, *grpcapi.ShortenRequest, ...grpc.CallOption) *grpcapi.ShortenResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpcapi.ShortenResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpcapi.ShortenRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}