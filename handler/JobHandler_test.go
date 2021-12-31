package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/mocks/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
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
