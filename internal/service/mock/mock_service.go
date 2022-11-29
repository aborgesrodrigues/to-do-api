// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/models.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	common "github.com/aborgesrodrigues/to-do-api/internal/common"
	gomock "github.com/golang/mock/gomock"
)

// MockSVCInterface is a mock of SVCInterface interface.
type MockSVCInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSVCInterfaceMockRecorder
}

// MockSVCInterfaceMockRecorder is the mock recorder for MockSVCInterface.
type MockSVCInterfaceMockRecorder struct {
	mock *MockSVCInterface
}

// NewMockSVCInterface creates a new mock instance.
func NewMockSVCInterface(ctrl *gomock.Controller) *MockSVCInterface {
	mock := &MockSVCInterface{ctrl: ctrl}
	mock.recorder = &MockSVCInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSVCInterface) EXPECT() *MockSVCInterfaceMockRecorder {
	return m.recorder
}

// AddTask mocks base method.
func (m *MockSVCInterface) AddTask(task *common.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTask", task)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTask indicates an expected call of AddTask.
func (mr *MockSVCInterfaceMockRecorder) AddTask(task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTask", reflect.TypeOf((*MockSVCInterface)(nil).AddTask), task)
}

// AddUser mocks base method.
func (m *MockSVCInterface) AddUser(user *common.User) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", user)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockSVCInterfaceMockRecorder) AddUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockSVCInterface)(nil).AddUser), user)
}

// DeleteTask mocks base method.
func (m *MockSVCInterface) DeleteTask(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockSVCInterfaceMockRecorder) DeleteTask(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockSVCInterface)(nil).DeleteTask), id)
}

// DeleteUser mocks base method.
func (m *MockSVCInterface) DeleteUser(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockSVCInterfaceMockRecorder) DeleteUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockSVCInterface)(nil).DeleteUser), id)
}

// GetTask mocks base method.
func (m *MockSVCInterface) GetTask(id string) (*common.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", id)
	ret0, _ := ret[0].(*common.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockSVCInterfaceMockRecorder) GetTask(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockSVCInterface)(nil).GetTask), id)
}

// GetUser mocks base method.
func (m *MockSVCInterface) GetUser(id string) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockSVCInterfaceMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockSVCInterface)(nil).GetUser), id)
}

// ListTasks mocks base method.
func (m *MockSVCInterface) ListTasks() ([]common.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasks")
	ret0, _ := ret[0].([]common.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks.
func (mr *MockSVCInterfaceMockRecorder) ListTasks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockSVCInterface)(nil).ListTasks))
}

// ListUserTasks mocks base method.
func (m *MockSVCInterface) ListUserTasks(id string) ([]common.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserTasks", id)
	ret0, _ := ret[0].([]common.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserTasks indicates an expected call of ListUserTasks.
func (mr *MockSVCInterfaceMockRecorder) ListUserTasks(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserTasks", reflect.TypeOf((*MockSVCInterface)(nil).ListUserTasks), id)
}

// ListUsers mocks base method.
func (m *MockSVCInterface) ListUsers() ([]common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers")
	ret0, _ := ret[0].([]common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockSVCInterfaceMockRecorder) ListUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockSVCInterface)(nil).ListUsers))
}

// UpdateTask mocks base method.
func (m *MockSVCInterface) UpdateTask(task *common.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", task)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockSVCInterfaceMockRecorder) UpdateTask(task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockSVCInterface)(nil).UpdateTask), task)
}

// UpdateUser mocks base method.
func (m *MockSVCInterface) UpdateUser(user *common.User) (*common.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(*common.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockSVCInterfaceMockRecorder) UpdateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockSVCInterface)(nil).UpdateUser), user)
}
