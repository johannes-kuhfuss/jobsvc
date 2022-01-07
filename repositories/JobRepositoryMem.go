package repositories

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type JobRepositoryMem struct {
	jobList map[string]domain.Job
	mu      *sync.Mutex
}

func NewJobRepositoryMem() JobRepositoryMem {
	jList := make(map[string]domain.Job)
	m := sync.Mutex{}
	return JobRepositoryMem{jList, &m}
}

func (jrm JobRepositoryMem) FindAll(status string) (*[]domain.Job, api_error.ApiErr) {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	if len(jrm.jobList) == 0 {
		msg := "No jobs in job list"
		logger.Info(msg)
		return nil, api_error.NewNotFoundError(msg)
	}
	if strings.TrimSpace(status) == "" {
		return convertMapToSlice(jrm.jobList), nil
	} else {
		return filterByStatus(jrm.jobList, status)
	}
}

func convertMapToSlice(jList map[string]domain.Job) *[]domain.Job {
	slice := make([]domain.Job, 0)
	for _, job := range jList {
		slice = append(slice, job)
	}
	return &slice
}

func filterByStatus(jList map[string]domain.Job, status string) (*[]domain.Job, api_error.ApiErr) {
	filteredByStatusList := make([]domain.Job, 0)
	for _, curJob := range jList {
		if curJob.Status == domain.JobStatus(status) {
			filteredByStatusList = append(filteredByStatusList, curJob)
		}
	}
	if len(filteredByStatusList) == 0 {
		msg := fmt.Sprintf("No jobs with status %v in joblist", status)
		logger.Info(msg)
		return nil, api_error.NewNotFoundError(msg)
	} else {
		return &filteredByStatusList, nil
	}
}

func (jrm JobRepositoryMem) FindById(id string) (*domain.Job, api_error.ApiErr) {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	if len(jrm.jobList) == 0 {
		msg := "No jobs in joblist"
		logger.Warn(msg)
		return nil, api_error.NewNotFoundError(msg)
	}
	return filterById(jrm.jobList, id)
}

func filterById(jList map[string]domain.Job, id string) (*domain.Job, api_error.ApiErr) {
	for _, curJob := range jList {
		if curJob.Id.String() == id {
			return &curJob, nil
		}
	}
	msg := fmt.Sprintf("No job with id %v in joblist", id)
	logger.Info(msg)
	return nil, api_error.NewNotFoundError(msg)
}

func (jrm JobRepositoryMem) Store(job domain.Job) api_error.ApiErr {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	job.ModifiedAt = date.GetNowUtc()
	jrm.jobList[job.Id.String()] = job
	return nil
}

func (jrm JobRepositoryMem) DeleteById(id string) api_error.ApiErr {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	if len(jrm.jobList) == 0 {
		msg := "No jobs in joblist"
		logger.Info(msg)
		return api_error.NewNotFoundError(msg)
	}
	_, err := filterById(jrm.jobList, id)
	if err != nil {
		return err
	}
	delete(jrm.jobList, id)
	return nil
}

func (jrm JobRepositoryMem) Dequeue(jobType string) (*domain.Job, api_error.ApiErr) {
	var nextJobId string = ""
	var nextJobDate time.Time = date.GetNowUtc().Add(1 * time.Second)

	jrm.mu.Lock()
	defer jrm.mu.Unlock()

	if len(jrm.jobList) == 0 {
		msg := "No jobs in joblist"
		logger.Info(msg)
		return nil, api_error.NewNotFoundError(msg)
	}
	for _, job := range jrm.jobList {
		if (job.Type == jobType) && (job.Status == domain.StatusCreated) {
			if job.CreatedAt.Before(nextJobDate) {
				nextJobDate = job.CreatedAt
				nextJobId = job.Id.String()
			}
		}
	}
	if nextJobId == "" {
		msg := "No more jobs to dequeue"
		logger.Info(msg)
		err := api_error.NewNotFoundError(msg)
		return nil, err
	}
	job, _ := filterById(jrm.jobList, nextJobId)
	return job, nil
}

func (jrm JobRepositoryMem) SetStatusById(id string, newStatus string, message string) api_error.ApiErr {
	job, err := jrm.FindById(id)
	if err != nil {
		return err
	}
	job.Status = domain.JobStatus(newStatus)
	job.AddHistory(message)
	jrm.Store(*job)
	return nil
}

func (jrm JobRepositoryMem) SetHistoryById(id string, message string) api_error.ApiErr {
	job, err := jrm.FindById(id)
	if err != nil {
		return err
	}
	job.AddHistory(message)
	jrm.Store(*job)
	return nil
}

func (jrm JobRepositoryMem) Update(id string, jobUpdReq dto.CreateUpdateJobRequest) (*domain.Job, api_error.ApiErr) {
	oldJob, err := jrm.FindById(id)
	if err != nil {
		return nil, err
	}
	updJob := mergeJobs(oldJob, jobUpdReq)
	jrm.Store(*updJob)
	return updJob, nil
}

func (jrm JobRepositoryMem) DeleteAllJobs() api_error.ApiErr {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	for key := range jrm.jobList {
		delete(jrm.jobList, key)
	}
	return nil
}
