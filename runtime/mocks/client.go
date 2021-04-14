// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	manifest "github.com/foldsh/fold/manifest"
	mock "github.com/stretchr/testify/mock"

	transport "github.com/foldsh/fold/runtime/transport"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// DoRequest provides a mock function with given fields: _a0, _a1
func (_m *Client) DoRequest(_a0 context.Context, _a1 *transport.Request) (*transport.Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *transport.Response
	if rf, ok := ret.Get(0).(func(context.Context, *transport.Request) *transport.Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transport.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *transport.Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManifest provides a mock function with given fields: _a0
func (_m *Client) GetManifest(_a0 context.Context) (*manifest.Manifest, error) {
	ret := _m.Called(_a0)

	var r0 *manifest.Manifest
	if rf, ok := ret.Get(0).(func(context.Context) *manifest.Manifest); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*manifest.Manifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Restart provides a mock function with given fields: _a0
func (_m *Client) Restart(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *Client) Start(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *Client) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}