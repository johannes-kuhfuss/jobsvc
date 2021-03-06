package service

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

//go:generate mockgen -destination=../mocks/service/mockJobService.go -package=service github.com/johannes-kuhfuss/jobsvc/service JobService
type JobService interface {
	CreateJob(dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr)
	GetAllJobs(dto.SortAndFilterRequest) (*[]dto.JobResponse, int, api_error.ApiErr)
	GetJobById(string) (*dto.JobResponse, api_error.ApiErr)
	DeleteJobById(string) api_error.ApiErr
	Dequeue(dto.DequeueRequest) (*dto.JobResponse, api_error.ApiErr)
	UpdateJob(string, dto.CreateUpdateJobRequest) (*dto.JobResponse, api_error.ApiErr)
	SetStatusById(string, dto.UpdateJobStatusRequest) api_error.ApiErr
	SetHistoryById(string, dto.UpdateJobHistoryRequest) api_error.ApiErr
	DeleteAllJobs() api_error.ApiErr
	CleanJobs() api_error.ApiErr
}

type DefaultJobService struct {
	repo domain.JobRepository
	Cfg  *config.AppConfig
}

func NewJobService(cfg *config.AppConfig, repository domain.JobRepository) DefaultJobService {
	return DefaultJobService{
		repo: repository,
		Cfg:  cfg,
	}
}

func (s DefaultJobService) GetAllJobs(safReq dto.SortAndFilterRequest) (*[]dto.JobResponse, int, api_error.ApiErr) {
	jobs, totalCount, err := s.repo.FindAll(safReq)
	if err != nil {
		return nil, 0, err
	}
	response := make([]dto.JobResponse, 0)
	for _, job := range *jobs {
		response = append(response, job.ToJobResponseDto())
	}
	return &response, totalCount, nil
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
	if strings.TrimSpace(statusReq.Message) == "" {
		message = fmt.Sprintf("Job status changed. New status: %v", statusReq.Status)
	} else {
		message = fmt.Sprintf("Job status changed. New status: %v; %v", statusReq.Status, statusReq.Message)
	}
	err = s.repo.SetStatusById(id, statusReq.Status, message)
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
	err = s.repo.SetHistoryById(id, historyReq.Message)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultJobService) DeleteAllJobs() api_error.ApiErr {
	err := s.repo.DeleteAllJobs()
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultJobService) CleanJobs() api_error.ApiErr {
	err := s.repo.CleanupJobs()
	if err != nil {
		return err
	}
	return nil
}
