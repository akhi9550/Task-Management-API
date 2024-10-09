// Code generated by MockGen. DO NOT EDIT.
// Source: pkg\usecase\interface\task.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	models "taskmanagementapi/pkg/utils/models"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskUseCase is a mock of TaskUseCase interface.
type MockTaskUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockTaskUseCaseMockRecorder
}

// MockTaskUseCaseMockRecorder is the mock recorder for MockTaskUseCase.
type MockTaskUseCaseMockRecorder struct {
	mock *MockTaskUseCase
}

// NewMockTaskUseCase creates a new mock instance.
func NewMockTaskUseCase(ctrl *gomock.Controller) *MockTaskUseCase {
	mock := &MockTaskUseCase{ctrl: ctrl}
	mock.recorder = &MockTaskUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskUseCase) EXPECT() *MockTaskUseCaseMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskUseCase) CreateTask(arg0 models.CreateTask, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskUseCaseMockRecorder) CreateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskUseCase)(nil).CreateTask), arg0, arg1)
}

// DeleteTask mocks base method.
func (m *MockTaskUseCase) DeleteTask(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskUseCaseMockRecorder) DeleteTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskUseCase)(nil).DeleteTask), arg0, arg1)
}

// GetTask mocks base method.
func (m *MockTaskUseCase) GetTask(arg0, arg1 string) (models.TaskDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", arg0, arg1)
	ret0, _ := ret[0].(models.TaskDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskUseCaseMockRecorder) GetTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskUseCase)(nil).GetTask), arg0, arg1)
}

// GetTasks mocks base method.
func (m *MockTaskUseCase) GetTasks(arg0 string) ([]models.TaskDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasks", arg0)
	ret0, _ := ret[0].([]models.TaskDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasks indicates an expected call of GetTasks.
func (mr *MockTaskUseCaseMockRecorder) GetTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasks", reflect.TypeOf((*MockTaskUseCase)(nil).GetTasks), arg0)
}

// UpdateTask mocks base method.
func (m *MockTaskUseCase) UpdateTask(arg0, arg1 string, arg2 models.CreateTask) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskUseCaseMockRecorder) UpdateTask(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskUseCase)(nil).UpdateTask), arg0, arg1, arg2)
}
