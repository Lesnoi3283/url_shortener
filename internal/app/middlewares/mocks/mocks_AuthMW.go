// Code generated by MockGen. DO NOT EDIT.
// Source: auth_mw.go

// Package mocks_MW is a generated GoMock package.
package mocks_MW

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserCreater is a mock of UserCreater interface.
type MockUserCreater struct {
	ctrl     *gomock.Controller
	recorder *MockUserCreaterMockRecorder
}

// MockUserCreaterMockRecorder is the mock recorder for MockUserCreater.
type MockUserCreaterMockRecorder struct {
	mock *MockUserCreater
}

// NewMockUserCreater creates a new mock instance.
func NewMockUserCreater(ctrl *gomock.Controller) *MockUserCreater {
	mock := &MockUserCreater{ctrl: ctrl}
	mock.recorder = &MockUserCreaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserCreater) EXPECT() *MockUserCreaterMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserCreater) CreateUser(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserCreaterMockRecorder) CreateUser(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserCreater)(nil).CreateUser), ctx)
}
