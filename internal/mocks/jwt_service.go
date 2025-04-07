// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/personal/task-management/pkg/utils/jwt (interfaces: JWTTokenServicer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	jwt "github.com/personal/task-management/pkg/utils/jwt"
)

// MockJWTTokenServicer is a mock of JWTTokenServicer interface.
type MockJWTTokenServicer struct {
	ctrl     *gomock.Controller
	recorder *MockJWTTokenServicerMockRecorder
}

// MockJWTTokenServicerMockRecorder is the mock recorder for MockJWTTokenServicer.
type MockJWTTokenServicerMockRecorder struct {
	mock *MockJWTTokenServicer
}

// NewMockJWTTokenServicer creates a new mock instance.
func NewMockJWTTokenServicer(ctrl *gomock.Controller) *MockJWTTokenServicer {
	mock := &MockJWTTokenServicer{ctrl: ctrl}
	mock.recorder = &MockJWTTokenServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTTokenServicer) EXPECT() *MockJWTTokenServicerMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJWTTokenServicer) GenerateToken(arg0 uuid.UUID, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJWTTokenServicerMockRecorder) GenerateToken(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJWTTokenServicer)(nil).GenerateToken), arg0, arg1, arg2)
}

// ValidateToken mocks base method.
func (m *MockJWTTokenServicer) ValidateToken(arg0 string) (*jwt.UserClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(*jwt.UserClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockJWTTokenServicerMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockJWTTokenServicer)(nil).ValidateToken), arg0)
}
