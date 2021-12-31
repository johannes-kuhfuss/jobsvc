// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/johannes-kuhfuss/jobsvc/service (interfaces: JobService)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/johannes-kuhfuss/jobsvc/dto"
	api_error "github.com/johannes-kuhfuss/services_utils/api_error"
)

// MockJobService is a mock of JobService interface.
type MockJobService struct {
	ctrl     *gomock.Controller
	recorder *MockJobServiceMockRecorder
}

// MockJobServiceMockRecorder is the mock recorder for MockJobService.
type MockJobServiceMockRecorder struct {
	mock *MockJobService
}

// NewMockJobService creates a new mock instance.
func NewMockJobService(ctrl *gomock.Controller) *MockJobService {
	mock := &MockJobService{ctrl: ctrl}
	mock.recorder = &MockJobServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobService) EXPECT() *MockJobServiceMockRecorder {
	return m.recorder
}

// CreateJob mocks base method.
func (m *MockJobService) CreateJob(arg0 dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJob", arg0)
	ret0, _ := ret[0].(*dto.JobResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// CreateJob indicates an expected call of CreateJob.
func (mr *MockJobServiceMockRecorder) CreateJob(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJob", reflect.TypeOf((*MockJobService)(nil).CreateJob), arg0)
}

// DeleteJobById mocks base method.
func (m *MockJobService) DeleteJobById(arg0 string) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteJobById", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// DeleteJobById indicates an expected call of DeleteJobById.
func (mr *MockJobServiceMockRecorder) DeleteJobById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteJobById", reflect.TypeOf((*MockJobService)(nil).DeleteJobById), arg0)
}

// Dequeue mocks base method.
func (m *MockJobService) Dequeue(arg0 dto.DequeueRequest) (*dto.JobResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dequeue", arg0)
	ret0, _ := ret[0].(*dto.JobResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// Dequeue indicates an expected call of Dequeue.
func (mr *MockJobServiceMockRecorder) Dequeue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dequeue", reflect.TypeOf((*MockJobService)(nil).Dequeue), arg0)
}

// GetAllJobs mocks base method.
func (m *MockJobService) GetAllJobs(arg0 string) (*[]dto.JobResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllJobs", arg0)
	ret0, _ := ret[0].(*[]dto.JobResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// GetAllJobs indicates an expected call of GetAllJobs.
func (mr *MockJobServiceMockRecorder) GetAllJobs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllJobs", reflect.TypeOf((*MockJobService)(nil).GetAllJobs), arg0)
}

// GetJobById mocks base method.
func (m *MockJobService) GetJobById(arg0 string) (*dto.JobResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJobById", arg0)
	ret0, _ := ret[0].(*dto.JobResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// GetJobById indicates an expected call of GetJobById.
func (mr *MockJobServiceMockRecorder) GetJobById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJobById", reflect.TypeOf((*MockJobService)(nil).GetJobById), arg0)
}

// SetHistoryById mocks base method.
func (m *MockJobService) SetHistoryById(arg0 string, arg1 dto.UpdateJobHistoryRequest) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHistoryById", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// SetHistoryById indicates an expected call of SetHistoryById.
func (mr *MockJobServiceMockRecorder) SetHistoryById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHistoryById", reflect.TypeOf((*MockJobService)(nil).SetHistoryById), arg0, arg1)
}

// SetStatusById mocks base method.
func (m *MockJobService) SetStatusById(arg0 string, arg1 dto.UpdateJobStatusRequest) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatusById", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// SetStatusById indicates an expected call of SetStatusById.
func (mr *MockJobServiceMockRecorder) SetStatusById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatusById", reflect.TypeOf((*MockJobService)(nil).SetStatusById), arg0, arg1)
}

// UpdateJob mocks base method.
func (m *MockJobService) UpdateJob(arg0 string, arg1 dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateJob", arg0, arg1)
	ret0, _ := ret[0].(*dto.JobResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// UpdateJob indicates an expected call of UpdateJob.
func (mr *MockJobServiceMockRecorder) UpdateJob(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateJob", reflect.TypeOf((*MockJobService)(nil).UpdateJob), arg0, arg1)
}
