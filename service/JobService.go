package service

import (
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobService interface {
	CreateJob(dto.CreateJobRequest) (*dto.JobResponse, api_error.ApiErr)
	GetAllJobs(string) (*[]dto.JobResponse, api_error.ApiErr)
	/*

		GetJobById(string) (*dto.JobResponse, api_error.ApiErr)
		CreateJob(dto.NewJobRequest) (*dto.JobResponse, api_error.ApiErr)
		DeleteJobById(string) api_error.ApiErr
		GetNextJob() (*dto.JobResponse, api_error.ApiErr)
		SetStatus(string, dto.JobStatusUpdateRequest) api_error.ApiErr
		SetResult(string, string) api_error.ApiErr
	*/
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

func (s DefaultJobService) CreateJob(jobreq dto.CreateJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	newJob, err := domain.NewJob(jobreq.Name, jobreq.Type)
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
