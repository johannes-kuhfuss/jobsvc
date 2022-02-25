// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/johannes-kuhfuss/jobsvc/domain (interfaces: JobRepository)

// Package domain is a generated GoMock package.
package domain

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/johannes-kuhfuss/jobsvc/domain"
	dto "github.com/johannes-kuhfuss/jobsvc/dto"
	api_error "github.com/johannes-kuhfuss/services_utils/api_error"
)

// MockJobRepository is a mock of JobRepository interface.
type MockJobRepository struct {
	ctrl     *gomock.Controller
	recorder *MockJobRepositoryMockRecorder
}

// MockJobRepositoryMockRecorder is the mock recorder for MockJobRepository.
type MockJobRepositoryMockRecorder struct {
	mock *MockJobRepository
}

// NewMockJobRepository creates a new mock instance.
func NewMockJobRepository(ctrl *gomock.Controller) *MockJobRepository {
	mock := &MockJobRepository{ctrl: ctrl}
	mock.recorder = &MockJobRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobRepository) EXPECT() *MockJobRepositoryMockRecorder {
	return m.recorder
}

// CleanupJobs mocks base method.
func (m *MockJobRepository) CleanupJobs(arg0, arg1 int) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanupJobs", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// CleanupJobs indicates an expected call of CleanupJobs.
func (mr *MockJobRepositoryMockRecorder) CleanupJobs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanupJobs", reflect.TypeOf((*MockJobRepository)(nil).CleanupJobs), arg0, arg1)
}

// DeleteAllJobs mocks base method.
func (m *MockJobRepository) DeleteAllJobs() api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllJobs")
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// DeleteAllJobs indicates an expected call of DeleteAllJobs.
func (mr *MockJobRepositoryMockRecorder) DeleteAllJobs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllJobs", reflect.TypeOf((*MockJobRepository)(nil).DeleteAllJobs))
}

// DeleteById mocks base method.
func (m *MockJobRepository) DeleteById(arg0 string) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockJobRepositoryMockRecorder) DeleteById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockJobRepository)(nil).DeleteById), arg0)
}

// Dequeue mocks base method.
func (m *MockJobRepository) Dequeue(arg0 string) (*domain.Job, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dequeue", arg0)
	ret0, _ := ret[0].(*domain.Job)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// Dequeue indicates an expected call of Dequeue.
func (mr *MockJobRepositoryMockRecorder) Dequeue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dequeue", reflect.TypeOf((*MockJobRepository)(nil).Dequeue), arg0)
}

// FindAll mocks base method.
func (m *MockJobRepository) FindAll(arg0 dto.SortAndFilterRequest) (*[]domain.Job, int, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].(*[]domain.Job)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(api_error.ApiErr)
	return ret0, ret1, ret2
}

// FindAll indicates an expected call of FindAll.
func (mr *MockJobRepositoryMockRecorder) FindAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockJobRepository)(nil).FindAll), arg0)
}

// FindById mocks base method.
func (m *MockJobRepository) FindById(arg0 string) (*domain.Job, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", arg0)
	ret0, _ := ret[0].(*domain.Job)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockJobRepositoryMockRecorder) FindById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockJobRepository)(nil).FindById), arg0)
}

// SetHistoryById mocks base method.
func (m *MockJobRepository) SetHistoryById(arg0, arg1 string) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHistoryById", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// SetHistoryById indicates an expected call of SetHistoryById.
func (mr *MockJobRepositoryMockRecorder) SetHistoryById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHistoryById", reflect.TypeOf((*MockJobRepository)(nil).SetHistoryById), arg0, arg1)
}

// SetStatusById mocks base method.
func (m *MockJobRepository) SetStatusById(arg0, arg1, arg2 string) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatusById", arg0, arg1, arg2)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// SetStatusById indicates an expected call of SetStatusById.
func (mr *MockJobRepositoryMockRecorder) SetStatusById(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatusById", reflect.TypeOf((*MockJobRepository)(nil).SetStatusById), arg0, arg1, arg2)
}

// Store mocks base method.
func (m *MockJobRepository) Store(arg0 domain.Job) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockJobRepositoryMockRecorder) Store(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockJobRepository)(nil).Store), arg0)
}

// Update mocks base method.
func (m *MockJobRepository) Update(arg0 string, arg1 dto.CreateUpdateJobRequest) (*domain.Job, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*domain.Job)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockJobRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockJobRepository)(nil).Update), arg0, arg1)
}
