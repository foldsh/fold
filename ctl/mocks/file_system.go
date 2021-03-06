// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// FileSystem is an autogenerated mock type for the FileSystem type
type FileSystem struct {
	mock.Mock
}

// MkdirAll provides a mock function with given fields: path, perm
func (_m *FileSystem) MkdirAll(path string, perm fs.FileMode) error {
	ret := _m.Called(path, perm)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, fs.FileMode) error); ok {
		r0 = rf(path, perm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
