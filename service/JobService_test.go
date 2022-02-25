package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/jobsvc/config"
	realdomain "github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/mocks/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	jobCtrl     *gomock.Controller
	mockJobRepo *domain.MockJobRepository
	jobService  JobService
	cfg         config.AppConfig
)

func setupJob(t *testing.T) func() {
	jobCtrl = gomock.NewController(t)
	mockJobRepo = domain.NewMockJobRepository(jobCtrl)
	jobService = NewJobService(&cfg, mockJobRepo)
	return func() {
		jobService = nil
		jobCtrl.Finish()
	}
}

func Test_GetAllJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("no jobs found")
	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}
	mockJobRepo.EXPECT().FindAll(safReq).Return(nil, 0, apiError)

	result, totalCount, err := jobService.GetAllJobs(safReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, 0, totalCount)
}

func Test_GetAllJobs_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	job1, _ := realdomain.NewJob("job 1", "encoding")
	job2, _ := realdomain.NewJob("job 2", "encoding")
	jobs := make([]realdomain.Job, 0)
	jobs = append(jobs, *job1)
	jobs = append(jobs, *job2)
	jobResult := make([]dto.JobResponse, 0)
	jobResult = append(jobResult, job1.ToJobResponseDto())
	jobResult = append(jobResult, job2.ToJobResponseDto())
	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}

	mockJobRepo.EXPECT().FindAll(safReq).Return(&jobs, len(jobs), nil)

	result, totalCount, err := jobService.GetAllJobs(safReq)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResult)
	assert.EqualValues(t, len(jobs), totalCount)
}

func Test_CreateJob_Returns_BaqRequestError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.CreateUpdateJobRequest{
		Name: "job 1",
	}
	result, err := jobService.CreateJob(jobReq)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Job must have a type", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_CreateJob_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.CreateUpdateJobRequest{
		Name: "job 1",
		Type: "encoding",
	}
	apiError := api_error.NewInternalServerError("database error", nil)
	mockJobRepo.EXPECT().Store(gomock.Any()).Return(apiError)

	result, err := jobService.CreateJob(jobReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_CreateJob_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.CreateUpdateJobRequest{
		Name: "job 1",
		Type: "encoding",
	}
	mockJobRepo.EXPECT().Store(gomock.Any()).Return(nil)

	result, err := jobService.CreateJob(jobReq)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.EqualValues(t, jobReq.Name, result.Name)
	assert.EqualValues(t, jobReq.Type, result.Type)
	assert.EqualValues(t, "created", result.Status)
}

func Test_GetJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	result, err := jobService.GetJobById(id)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_GetJobById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	jobResp := newJob.ToJobResponseDto()
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)

	result, err := jobService.GetJobById(id)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResp)
}

func Test_DeleteJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_DeleteJobById_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	apiError := api_error.NewInternalServerError("database error", nil)
	mockJobRepo.EXPECT().DeleteById(id).Return(apiError)

	err := jobService.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_DeleteJobById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().DeleteById(id).Return(nil)

	err := jobService.DeleteJobById(id)

	assert.Nil(t, err)
}

func Test_Dequeue_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("No next job found")
	dqReq := dto.DequeueRequest{
		Type: "Encoding",
	}
	mockJobRepo.EXPECT().Dequeue(dqReq.Type).Return(nil, apiError)

	job, err := jobService.Dequeue(dqReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_Dequeue_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	nextJob, _ := realdomain.NewJob("job 1", "encoding")
	dqReq := dto.DequeueRequest{
		Type: "Encoding",
	}
	mockJobRepo.EXPECT().Dequeue(dqReq.Type).Return(nextJob, nil)

	job, err := jobService.Dequeue(dqReq)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, nextJob.Name, job.Name)
	assert.EqualValues(t, nextJob.Type, job.Type)
}

func Test_UpdateJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	updReq := dto.CreateUpdateJobRequest{}
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	job, err := jobService.UpdateJob(id, updReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_UpdateJob_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.CreateUpdateJobRequest{}
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	apiError := api_error.NewInternalServerError("database error", nil)
	mockJobRepo.EXPECT().Update(id, updReq).Return(nil, apiError)

	job, err := jobService.UpdateJob(id, updReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_UpdateJob_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.CreateUpdateJobRequest{}
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().Update(id, updReq).Return(newJob, nil)

	job, err := jobService.UpdateJob(id, updReq)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, newJob.Name, job.Name)
	assert.EqualValues(t, newJob.Type, job.Type)
}

func Test_SetStatusById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	updReq := dto.UpdateJobStatusRequest{
		Status:  "",
		Message: "",
	}
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.SetStatusById(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_SetStatusById_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.UpdateJobStatusRequest{
		Status:  "running",
		Message: "",
	}
	msg := fmt.Sprintf("Job status changed. New status: %v", updReq.Status)
	apiError := api_error.NewInternalServerError("Database error", nil)
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetStatusById(id, updReq.Status, msg).Return(apiError)

	err := jobService.SetStatusById(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_SetStatusById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.UpdateJobStatusRequest{
		Status:  "running",
		Message: "oops",
	}
	msg := fmt.Sprintf("Job status changed. New status: %v; %v", updReq.Status, updReq.Message)
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetStatusById(id, updReq.Status, msg).Return(nil)

	err := jobService.SetStatusById(id, updReq)

	assert.Nil(t, err)
}

func Test_SetHistoryById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	updReq := dto.UpdateJobHistoryRequest{
		Message: "",
	}
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.SetHistoryById(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_SetHistoryById_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.UpdateJobHistoryRequest{
		Message: "new message",
	}
	apiError := api_error.NewInternalServerError("Database error", nil)
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetHistoryById(id, updReq.Message).Return(apiError)

	err := jobService.SetHistoryById(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_SetHistoryById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "encoding")
	id := newJob.Id.String()
	updReq := dto.UpdateJobHistoryRequest{
		Message: "new message",
	}
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetHistoryById(id, updReq.Message).Return(nil)

	err := jobService.SetHistoryById(id, updReq)

	assert.Nil(t, err)
}

func Test_DeleteAllJobs_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("Database error", nil)
	mockJobRepo.EXPECT().DeleteAllJobs().Return(apiError)

	err := jobService.DeleteAllJobs()

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, apiError.StatusCode(), err.StatusCode())
}

func Test_DeleteAllJobs_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	mockJobRepo.EXPECT().DeleteAllJobs().Return(nil)

	err := jobService.DeleteAllJobs()

	assert.Nil(t, err)
}
