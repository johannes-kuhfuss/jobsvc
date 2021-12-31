package service

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobService interface {
	CreateJob(dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr)
	GetAllJobs(string) (*[]dto.JobResponse, api_error.ApiErr)
	GetJobById(string) (*dto.JobResponse, api_error.ApiErr)
	DeleteJobById(string) api_error.ApiErr
	Dequeue(dto.DequeueRequest) (*dto.JobResponse, api_error.ApiErr)
	UpdateJob(string, dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr)
	SetStatusById(string, dto.UpdateJobStatusRequest) api_error.ApiErr
	SetHistoryById(string, dto.UpdateJobHistoryRequest) api_error.ApiErr
}

type DefaultJobService struct {
	repo domain.JobRepository
}

func NewJobService(repository domain.JobRepository) DefaultJobService {
	return DefaultJobService{repository}
}

func (s DefaultJobService) GetAllJobs(status string) (*[]dto.JobResponse, api_error.ApiErr) {
	jobs, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	response := make([]dto.JobResponse, 0)
	for _, job := range *jobs {
		response = append(response, job.ToJobResponseDto())
	}
	return &response, nil
}

func (s DefaultJobService) CreateJob(jobReq dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	newJob, err := domain.NewJobFromJobRequestDto(jobReq)
	if err != nil {
		return nil, err
	}
	err = s.repo.Store(*newJob)
	if err != nil {
		return nil, err
	}
	response := newJob.ToJobResponseDto()
	return &response, nil
}

func (s DefaultJobService) GetJobById(id string) (*dto.JobResponse, api_error.ApiErr) {
	job, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	response := job.ToJobResponseDto()
	return &response, nil
}

func (s DefaultJobService) DeleteJobById(id string) api_error.ApiErr {
	_, err := s.GetJobById(id)
	if err != nil {
		return api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	}
	err = s.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultJobService) Dequeue(dqReq dto.DequeueRequest) (*dto.JobResponse, api_error.ApiErr) {
	job, err := s.repo.Dequeue(dqReq.Type)
	if err != nil {
		return nil, err
	}
	response := job.ToJobResponseDto()
	return &response, nil
}

func (s DefaultJobService) UpdateJob(id string, jobReq dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	_, err := s.GetJobById(id)
	if err != nil {
		return nil, api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	}
	newJob, err := s.repo.Update(id, jobReq)
	if err != nil {
		return nil, err
	}
	response := newJob.ToJobResponseDto()
	return &response, nil
}

func (s DefaultJobService) SetStatusById(id string, statusReq dto.UpdateJobStatusRequest) api_error.ApiErr {
	var message string
	_, err := s.GetJobById(id)
	if err != nil {
		return api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	}
	statusVal := string(domain.JobStatus(statusReq.Status))
	if strings.TrimSpace(statusVal) == "" {
		return api_error.NewBadRequestError(fmt.Sprintf("Wrong status value %v when updating job status", statusVal))
	}
	if strings.TrimSpace(statusReq.Message) == "" {
		message = fmt.Sprintf("Job status changed. New status: %v", statusVal)
	} else {
		message = fmt.Sprintf("Job status changed. New status: %v; %v", statusVal, statusReq.Message)
	}
	err = s.repo.SetStatusById(id, statusVal, message)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultJobService) SetHistoryById(id string, historyReq dto.UpdateJobHistoryRequest) api_error.ApiErr {
	_, err := s.GetJobById(id)
	if err != nil {
		return api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	}
	if strings.TrimSpace(historyReq.Message) == "" {
		return api_error.NewBadRequestError("Empty message when updating job history")
	}
	err = s.repo.SetHistoryById(id, historyReq.Message)
	if err != nil {
		return err
	}
	return nil
}
