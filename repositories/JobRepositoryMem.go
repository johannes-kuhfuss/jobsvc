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
		return nil, api_error.NewNotFoundError("no jobs in joblist")
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
		return nil, api_error.NewNotFoundError(fmt.Sprintf("no jobs with status %v in joblist", status))
	} else {
		return &filteredByStatusList, nil
	}
}

func (jrm JobRepositoryMem) FindById(id string) (*domain.Job, api_error.ApiErr) {
	jrm.mu.Lock()
	defer jrm.mu.Unlock()
	if len(jrm.jobList) == 0 {
		return nil, api_error.NewNotFoundError("no jobs in joblist")
	}
	return filterById(jrm.jobList, id)
}

func filterById(jList map[string]domain.Job, id string) (*domain.Job, api_error.ApiErr) {
	for _, curJob := range jList {
		if curJob.Id.String() == id {
			return &curJob, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("no job with id %v in joblist", id))
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
		return api_error.NewNotFoundError("no jobs in joblist")
	}
	_, err := filterById(jrm.jobList, id)
	if err != nil {
		return err
	}
	delete(jrm.jobList, id)
	return nil
}

func (jrm JobRepositoryMem) GetNext() (*domain.Job, api_error.ApiErr) {
	var nextJobId string = ""
	var nextJobDate time.Time = date.GetNowUtc().Add(1 * time.Second)

	jrm.mu.Lock()
	defer jrm.mu.Unlock()

	if len(jrm.jobList) == 0 {
		err := api_error.NewNotFoundError("no jobs in joblist")
		return nil, err
	}
	for _, job := range jrm.jobList {
		if job.Status == domain.StatusCreated {
			if job.CreatedAt.Before(nextJobDate) {
				nextJobDate = job.CreatedAt
				nextJobId = job.Id.String()
			}
		}
	}
	if nextJobId == "" {
		err := api_error.NewNotFoundError("no jobs with status created in joblist")
		return nil, err
	}
	job, _ := filterById(jrm.jobList, nextJobId)
	return job, nil
}

func (jrm JobRepositoryMem) SetStatus(id string, newStatus dto.UpdateJobStatusRequest) api_error.ApiErr {
	job, err := jrm.FindById(id)
	if err != nil {
		return err
	}
	job.Status = domain.JobStatus(newStatus.Status)
	jrm.Store(*job)
	return nil
}
