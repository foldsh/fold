// Code generated by MockGen. DO NOT EDIT.
// Source: runtime/service.go

// Package runtime_test is a generated GoMock package.
package runtime_test

import (
	context "context"
	os "os"
	reflect "reflect"

	types "github.com/foldsh/fold/runtime/types"
	gomock "github.com/golang/mock/gomock"
)

// MockSupervisor is a mock of Supervisor interface.
type MockSupervisor struct {
	ctrl     *gomock.Controller
	recorder *MockSupervisorMockRecorder
}

// MockSupervisorMockRecorder is the mock recorder for MockSupervisor.
type MockSupervisorMockRecorder struct {
	mock *MockSupervisor
}

// NewMockSupervisor creates a new mock instance.
func NewMockSupervisor(ctrl *gomock.Controller) *MockSupervisor {
	mock := &MockSupervisor{ctrl: ctrl}
	mock.recorder = &MockSupervisorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSupervisor) EXPECT() *MockSupervisorMockRecorder {
	return m.recorder
}

// Kill mocks base method.
func (m *MockSupervisor) Kill() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Kill")
	ret0, _ := ret[0].(error)
	return ret0
}

// Kill indicates an expected call of Kill.
func (mr *MockSupervisorMockRecorder) Kill() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockSupervisor)(nil).Kill))
}

// Restart mocks base method.
func (m *MockSupervisor) Restart() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restart")
	ret0, _ := ret[0].(error)
	return ret0
}

// Restart indicates an expected call of Restart.
func (mr *MockSupervisorMockRecorder) Restart() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restart", reflect.TypeOf((*MockSupervisor)(nil).Restart))
}

// Signal mocks base method.
func (m *MockSupervisor) Signal(sig os.Signal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signal", sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signal indicates an expected call of Signal.
func (mr *MockSupervisorMockRecorder) Signal(sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*MockSupervisor)(nil).Signal), sig)
}

// Start mocks base method.
func (m *MockSupervisor) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockSupervisorMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockSupervisor)(nil).Start))
}

// Stop mocks base method.
func (m *MockSupervisor) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockSupervisorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockSupervisor)(nil).Stop))
}

// Wait mocks base method.
func (m *MockSupervisor) Wait() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockSupervisorMockRecorder) Wait() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockSupervisor)(nil).Wait))
}

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// DoRequest mocks base method.
func (m *MockClient) DoRequest(ctx context.Context, req *types.Request) (*types.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoRequest", ctx, req)
	ret0, _ := ret[0].(*types.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoRequest indicates an expected call of DoRequest.
func (mr *MockClientMockRecorder) DoRequest(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoRequest", reflect.TypeOf((*MockClient)(nil).DoRequest), ctx, req)
}

// GetManifest mocks base method.
func (m *MockClient) GetManifest(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManifest", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetManifest indicates an expected call of GetManifest.
func (mr *MockClientMockRecorder) GetManifest(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManifest", reflect.TypeOf((*MockClient)(nil).GetManifest), ctx)
}

// Start mocks base method.
func (m *MockClient) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockClientMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockClient)(nil).Start))
}
