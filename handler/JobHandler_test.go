package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/mocks/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/microcosm-cc/bluemonday"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	cfg         config.AppConfig
	jh          JobHandlers
	router      *gin.Engine
	mockService *service.MockJobService
	recorder    *httptest.ResponseRecorder
)

func setupTest(t *testing.T) func() {
	cfg.RunTime.BmPolicy = bluemonday.UGCPolicy()
	sani, _ := sanitize.New()
	cfg.RunTime.Sani = sani
	ctrl := gomock.NewController(t)
	mockService = service.NewMockJobService(ctrl)
	jh = JobHandlers{
		Service: mockService,
		Cfg:     &cfg,
	}
	jh.Cfg = &cfg
	router = gin.Default()
	recorder = httptest.NewRecorder()
	return func() {
		router = nil
		ctrl.Finish()
	}
}

func Test_getJobId_NonKsuid_Returns_BadRequestError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	testParam := "wrong_id"

	jobId, err := jh.getJobId(testParam)

	assert.NotNil(t, err)
	assert.EqualValues(t, "", jobId)
	assert.EqualValues(t, "User id should be a ksuid", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_getJobId_WithKsuid_Returns_String(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	testParam := ksuid.New()

	jobId, err := jh.getJobId(testParam.String())

	assert.NotNil(t, jobId)
	assert.Nil(t, err)
	assert.EqualValues(t, testParam.String(), jobId)
}

func Test_CreateJob_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("invalid json body for job creation")
	errorJson, _ := json.Marshal(apiError)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_CreateJob_Returns_InvalidInputError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("could not validate input data for create")
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.CreateUpdateJobRequest{
		Name:     "Job 1",
		Type:     "Encoding",
		Priority: "bogus",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_CreateJob_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.CreateUpdateJobRequest{
		Name: "Job 1",
		Type: "Encoding",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().CreateJob(jobReq).Return(nil, apiError)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_CreateJob_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	jobReq := dto.CreateUpdateJobRequest{
		Name: "Job 1",
		Type: "Encoding",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	jobResp := dto.JobResponse{
		Id:            ksuid.New().String(),
		CorrelationId: "",
		Name:          jobReq.Name,
		CreatedAt:     date.GetNowUtc(),
		CreatedBy:     "",
		ModifiedAt:    date.GetNowUtc(),
		ModifiedBy:    "",
		Status:        "created",
		Source:        "",
		Destination:   "",
		Type:          jobReq.Type,
		SubType:       "",
		Action:        "",
		ActionDetails: "",
		Progress:      0,
		History:       "",
		ExtraData:     "",
		Priority:      "medium",
		Rank:          0,
	}
	bodyJson, _ := json.Marshal(jobResp)
	mockService.EXPECT().CreateJob(jobReq).Return(&jobResp, nil)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusCreated, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}

func Test_GetAllJobs_Returns_BadRequestError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("database error")
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().GetAllJobs("").Return(nil, apiError)

	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_GetAllJobs_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	dummyJobList := createDummyJobList()
	dummyJobListJson, _ := json.Marshal(dummyJobList)
	mockService.EXPECT().GetAllJobs("").Return(&dummyJobList, nil)

	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, dummyJobListJson, recorder.Body.String())
}

func createDummyJobList() []dto.JobResponse {
	job1, _ := domain.NewJob("Job 1", "Encoding")
	job2, _ := domain.NewJob("Job 2", "Encondig")
	job1Dto := job1.ToJobResponseDto()
	job2Dto := job2.ToJobResponseDto()
	dummyJobList := []dto.JobResponse{}
	dummyJobList = append(dummyJobList, job1Dto)
	dummyJobList = append(dummyJobList, job2Dto)
	return dummyJobList
}

func Test_GetJobById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, "/jobs/not_a_ksuid", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_GetJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().GetJobById(gomock.Eq(id.String())).Return(nil, apiError)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNotFound, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_GetJobById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	newJob, _ := domain.NewJob("Job 1", "Encoding")
	newReq := newJob.ToJobResponseDto()
	bodyJson, _ := json.Marshal(newReq)
	mockService.EXPECT().GetJobById(id.String()).Return(&newReq, nil)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, "/jobs/not_a_ksuid", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().DeleteJobById(id.String()).Return(apiError)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNotFound, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	mockService.EXPECT().DeleteJobById(id.String()).Return(nil)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNoContent, recorder.Code)
}

func Test_Dequeue_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("invalid json body for dequeue request")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("/jobs/dequeue", jh.Dequeue)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/dequeue", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_Dequeue_Returns_InvalidInputError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("could not validate input data for dequeue")
	errorJson, _ := json.Marshal(apiError)
	req := dto.DequeueRequest{
		Type: "",
	}
	bodyJson, _ := json.Marshal(req)
	router.PUT("/jobs/dequeue", jh.Dequeue)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/dequeue", strings.NewReader(string(bodyJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_Dequeue_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	req := dto.DequeueRequest{
		Type: "Encoding",
	}
	bodyJson, _ := json.Marshal(req)
	router.PUT("/jobs/dequeue", jh.Dequeue)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/dequeue", strings.NewReader(string(bodyJson)))
	mockService.EXPECT().Dequeue(req).Return(nil, apiError)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_Dequeue_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	req := dto.DequeueRequest{
		Type: "Encoding",
	}
	bodyJson, _ := json.Marshal(req)
	newJob, _ := domain.NewJob("Job 1", "Encoding")
	jobResp := newJob.ToJobResponseDto()
	respJson, _ := json.Marshal(jobResp)
	router.PUT("/jobs/dequeue", jh.Dequeue)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/dequeue", strings.NewReader(string(bodyJson)))
	mockService.EXPECT().Dequeue(req).Return(&jobResp, nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, respJson, recorder.Body.String())
}

func Test_UpdateJob_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("/jobs/:job_id", jh.UpdateJob)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/not_a_ksuid", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_UpdateJob_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("invalid json body for job update")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("/jobs/:job_id", jh.UpdateJob)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_UpdateJob_Returns_InvalidInputError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("could not validate input data for update")
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.CreateUpdateJobRequest{
		Name:     "Job 1",
		Type:     "Encoding",
		Priority: "bogus",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	router.PUT("/jobs/:job_id", jh.UpdateJob)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_UpdateJob_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.CreateUpdateJobRequest{
		Priority: "high",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().UpdateJob(id.String(), jobReq).Return(nil, apiError)
	router.PUT("/jobs/:job_id", jh.UpdateJob)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_UpdateJob_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	jobReq := dto.CreateUpdateJobRequest{
		Priority: "high",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	newJob, _ := domain.NewJob("Job 1", "Encoding")
	newJobResp := newJob.ToJobResponseDto()
	newJobRespJson, _ := json.Marshal(newJobResp)
	mockService.EXPECT().UpdateJob(id.String(), jobReq).Return(&newJobResp, nil)
	router.PUT("/jobs/:job_id", jh.UpdateJob)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, newJobRespJson, recorder.Body.String())
}

func Test_SetStatusById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("jobs/:job_id/status", jh.SetStatusById)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/not_a_ksuid/status", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetStatusById_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("invalid json body for job status update")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("jobs/:job_id/status", jh.SetStatusById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/status", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetStatusById_Returns_InvalidInputError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("could not validate input data for update status")
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.UpdateJobStatusRequest{
		Status: "bogus",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	router.PUT("jobs/:job_id/status", jh.SetStatusById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/status", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetStatusById_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.UpdateJobStatusRequest{
		Status: "running",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().SetStatusById(id.String(), jobReq).Return(apiError)
	router.PUT("jobs/:job_id/status", jh.SetStatusById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/status", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetStatusById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	jobReq := dto.UpdateJobStatusRequest{
		Status: "running",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().SetStatusById(id.String(), jobReq).Return(nil)
	router.PUT("jobs/:job_id/status", jh.SetStatusById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/status", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNoContent, recorder.Code)
}

func Test_SetHistoryById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("jobs/:job_id/history", jh.SetHistoryById)
	request, _ := http.NewRequest(http.MethodPut, "/jobs/not_a_ksuid/history", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetHistoryById_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("invalid json body for job history update")
	errorJson, _ := json.Marshal(apiError)
	router.PUT("jobs/:job_id/history", jh.SetHistoryById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/history", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetHistoryById_Returns_InvalidInputError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewBadRequestError("could not validate input data for update history")
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.UpdateJobHistoryRequest{}
	jobReqJson, _ := json.Marshal(jobReq)
	router.PUT("jobs/:job_id/history", jh.SetHistoryById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/history", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetHistoryById_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.UpdateJobHistoryRequest{
		Message: "my message",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().SetHistoryById(id.String(), jobReq).Return(apiError)
	router.PUT("jobs/:job_id/history", jh.SetHistoryById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/history", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_SetHistoryById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	jobReq := dto.UpdateJobHistoryRequest{
		Message: "my message",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().SetHistoryById(id.String(), jobReq).Return(nil)
	router.PUT("jobs/:job_id/history", jh.SetHistoryById)
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/jobs/%v/history", id), strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNoContent, recorder.Code)
}
