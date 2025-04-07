// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/personal/task-management/internal/usecase (interfaces: WebSocketService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
	domain "github.com/personal/task-management/internal/domain"
)

// MockWebSocketService is a mock of WebSocketService interface.
type MockWebSocketService struct {
	ctrl     *gomock.Controller
	recorder *MockWebSocketServiceMockRecorder
}

// MockWebSocketServiceMockRecorder is the mock recorder for MockWebSocketService.
type MockWebSocketServiceMockRecorder struct {
	mock *MockWebSocketService
}

// NewMockWebSocketService creates a new mock instance.
func NewMockWebSocketService(ctrl *gomock.Controller) *MockWebSocketService {
	mock := &MockWebSocketService{ctrl: ctrl}
	mock.recorder = &MockWebSocketServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebSocketService) EXPECT() *MockWebSocketServiceMockRecorder {
	return m.recorder
}

// BroadcastTaskUpdate mocks base method.
func (m *MockWebSocketService) BroadcastTaskUpdate(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastTaskUpdate", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadcastTaskUpdate indicates an expected call of BroadcastTaskUpdate.
func (mr *MockWebSocketServiceMockRecorder) BroadcastTaskUpdate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastTaskUpdate", reflect.TypeOf((*MockWebSocketService)(nil).BroadcastTaskUpdate), arg0, arg1, arg2)
}

// CreateDirectRoom mocks base method.
func (m *MockWebSocketService) CreateDirectRoom(arg0, arg1 string) (*domain.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDirectRoom", arg0, arg1)
	ret0, _ := ret[0].(*domain.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDirectRoom indicates an expected call of CreateDirectRoom.
func (mr *MockWebSocketServiceMockRecorder) CreateDirectRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDirectRoom", reflect.TypeOf((*MockWebSocketService)(nil).CreateDirectRoom), arg0, arg1)
}

// CreateGroupRoom mocks base method.
func (m *MockWebSocketService) CreateGroupRoom(arg0 string, arg1 []string) (*domain.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroupRoom", arg0, arg1)
	ret0, _ := ret[0].(*domain.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroupRoom indicates an expected call of CreateGroupRoom.
func (mr *MockWebSocketServiceMockRecorder) CreateGroupRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroupRoom", reflect.TypeOf((*MockWebSocketService)(nil).CreateGroupRoom), arg0, arg1)
}

// HandleConnection mocks base method.
func (m *MockWebSocketService) HandleConnection(arg0 *websocket.Conn, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleConnection", arg0, arg1)
}

// HandleConnection indicates an expected call of HandleConnection.
func (mr *MockWebSocketServiceMockRecorder) HandleConnection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleConnection", reflect.TypeOf((*MockWebSocketService)(nil).HandleConnection), arg0, arg1)
}

// JoinRoom mocks base method.
func (m *MockWebSocketService) JoinRoom(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinRoom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinRoom indicates an expected call of JoinRoom.
func (mr *MockWebSocketServiceMockRecorder) JoinRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinRoom", reflect.TypeOf((*MockWebSocketService)(nil).JoinRoom), arg0, arg1)
}

// LeaveRoom mocks base method.
func (m *MockWebSocketService) LeaveRoom(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveRoom", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveRoom indicates an expected call of LeaveRoom.
func (mr *MockWebSocketServiceMockRecorder) LeaveRoom(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveRoom", reflect.TypeOf((*MockWebSocketService)(nil).LeaveRoom), arg0, arg1)
}

// SendDirectMessage mocks base method.
func (m *MockWebSocketService) SendDirectMessage(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDirectMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendDirectMessage indicates an expected call of SendDirectMessage.
func (mr *MockWebSocketServiceMockRecorder) SendDirectMessage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDirectMessage", reflect.TypeOf((*MockWebSocketService)(nil).SendDirectMessage), arg0, arg1, arg2)
}

// SendGroupMessage mocks base method.
func (m *MockWebSocketService) SendGroupMessage(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendGroupMessage", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendGroupMessage indicates an expected call of SendGroupMessage.
func (mr *MockWebSocketServiceMockRecorder) SendGroupMessage(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendGroupMessage", reflect.TypeOf((*MockWebSocketService)(nil).SendGroupMessage), arg0, arg1, arg2)
}
