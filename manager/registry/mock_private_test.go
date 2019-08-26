// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chmking/horde/protobuf/private (interfaces: AgentClient)

// Package registry_test is a generated GoMock package.
package registry_test

import (
	context "context"
	private "github.com/chmking/horde/protobuf/private"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockAgentClient is a mock of AgentClient interface
type MockAgentClient struct {
	ctrl     *gomock.Controller
	recorder *MockAgentClientMockRecorder
}

// MockAgentClientMockRecorder is the mock recorder for MockAgentClient
type MockAgentClientMockRecorder struct {
	mock *MockAgentClient
}

// NewMockAgentClient creates a new mock instance
func NewMockAgentClient(ctrl *gomock.Controller) *MockAgentClient {
	mock := &MockAgentClient{ctrl: ctrl}
	mock.recorder = &MockAgentClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAgentClient) EXPECT() *MockAgentClientMockRecorder {
	return m.recorder
}

// Healthcheck mocks base method
func (m *MockAgentClient) Healthcheck(arg0 context.Context, arg1 *private.HealthcheckRequest, arg2 ...grpc.CallOption) (*private.HealthcheckResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Healthcheck", varargs...)
	ret0, _ := ret[0].(*private.HealthcheckResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Healthcheck indicates an expected call of Healthcheck
func (mr *MockAgentClientMockRecorder) Healthcheck(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Healthcheck", reflect.TypeOf((*MockAgentClient)(nil).Healthcheck), varargs...)
}

// Quit mocks base method
func (m *MockAgentClient) Quit(arg0 context.Context, arg1 *private.QuitRequest, arg2 ...grpc.CallOption) (*private.QuitResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Quit", varargs...)
	ret0, _ := ret[0].(*private.QuitResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Quit indicates an expected call of Quit
func (mr *MockAgentClientMockRecorder) Quit(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Quit", reflect.TypeOf((*MockAgentClient)(nil).Quit), varargs...)
}

// Scale mocks base method
func (m *MockAgentClient) Scale(arg0 context.Context, arg1 *private.ScaleRequest, arg2 ...grpc.CallOption) (*private.ScaleResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scale", varargs...)
	ret0, _ := ret[0].(*private.ScaleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scale indicates an expected call of Scale
func (mr *MockAgentClientMockRecorder) Scale(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scale", reflect.TypeOf((*MockAgentClient)(nil).Scale), varargs...)
}

// Stop mocks base method
func (m *MockAgentClient) Stop(arg0 context.Context, arg1 *private.StopRequest, arg2 ...grpc.CallOption) (*private.StopResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stop", varargs...)
	ret0, _ := ret[0].(*private.StopResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop
func (mr *MockAgentClientMockRecorder) Stop(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAgentClient)(nil).Stop), varargs...)
}
