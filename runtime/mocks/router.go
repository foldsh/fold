// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	http "net/http"

	manifest "github.com/foldsh/fold/manifest"
	mock "github.com/stretchr/testify/mock"
)

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// Configure provides a mock function with given fields: _a0
func (_m *Router) Configure(_a0 *manifest.Manifest) {
	_m.Called(_a0)
}

// ServeHTTP provides a mock function with given fields: _a0, _a1
func (_m *Router) ServeHTTP(_a0 http.ResponseWriter, _a1 *http.Request) {
	_m.Called(_a0, _a1)
}
