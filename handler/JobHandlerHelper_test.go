package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateCreateUpdateJobRequest_NoType_Returns_BadRequestError(t *testing.T) {
	req := dto.CreateUpdateJobRequest{}

	err := validateCreateJobRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Job create / update request must have a type", err.Message())
}

func Test_validateCreateJobRequest_InvalidPriority_Returns_BadRequestError(t *testing.T) {
	prio := "bogus"
	req := dto.CreateUpdateJobRequest{
		Type:     "encoding",
		Priority: prio,
	}

	err := validateCreateJobRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("Priority value %v does not exist", prio), err.Message())
}

func Test_validateCreateJobRequest_ValidRequest_Returns_NoError(t *testing.T) {
	req := dto.CreateUpdateJobRequest{
		Type: "encoding",
	}

	err := validateCreateJobRequest(req)

	assert.Nil(t, err)
}

func Test_validateUpdateJobRequest_InvalidPriority_Returns_BadRequestError(t *testing.T) {
	prio := "bogus"
	req := dto.CreateUpdateJobRequest{
		Type:     "encoding",
		Priority: prio,
	}

	err := validateUpdateJobRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("Priority value %v does not exist", prio), err.Message())
}

func Test_validateUpdateJobRequest_ValidRequest_Returns_NoError(t *testing.T) {
	req := dto.CreateUpdateJobRequest{
		Type: "encoding",
	}

	err := validateUpdateJobRequest(req)

	assert.Nil(t, err)
}

func Test_validateDequeueRequest_NoType_Returns_BadRequestError(t *testing.T) {
	req := dto.DequeueRequest{}

	err := validateDequeueRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Dequeue request must have a type", err.Message())
}

func Test_validateDequeueRequest_ValidRequest_Returns_NoError(t *testing.T) {
	req := dto.DequeueRequest{
		Type: "encoding",
	}

	err := validateDequeueRequest(req)

	assert.Nil(t, err)
}

func Test_validateUpdateJobStatusRequest_NoStatus_Returns_BadRequestError(t *testing.T) {
	req := dto.UpdateJobStatusRequest{}

	err := validateUpdateJobStatusRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Update status request must have a status", err.Message())
}

func Test_validateUpdateJobStatusRequest_InvalidStatus_Returns_BadRequestError(t *testing.T) {

	req := dto.UpdateJobStatusRequest{
		Status: "bogus",
	}

	err := validateUpdateJobStatusRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("Wrong status value %v when updating job status", req.Status), err.Message())
}

func Test_validateUpdateJobStatusRequest_ValidRequest_Returns_NoError(t *testing.T) {
	req := dto.UpdateJobStatusRequest{
		Status: "paused",
	}

	err := validateUpdateJobStatusRequest(req)

	assert.Nil(t, err)
}

func Test_validateUpdateJobHistoryRequest_NoMessage_Returns_BadRequestError(t *testing.T) {
	req := dto.UpdateJobHistoryRequest{}

	err := validateUpdateJobHistoryRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Update history request must have a message", err.Message())
}

func Test_validateUpdateJobHistoryRequest_Returns_NoError(t *testing.T) {
	req := dto.UpdateJobHistoryRequest{
		Message: "Hey!",
	}

	err := validateUpdateJobHistoryRequest(req)

	assert.Nil(t, err)
}

func Test_extractSorts_NoInput_Returns_DefaultSort(t *testing.T) {
	var safParams url.Values

	sorts, err := extractSorts(safParams)

	assert.NotNil(t, sorts)
	assert.Nil(t, err)
	assert.EqualValues(t, "id", sorts[0].Field)
	assert.EqualValues(t, "DESC", sorts[0].Dir)
}

func Test_extractSorts_MalformedParam_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=asdf")
	safParams := url.Query()

	sorts, err := extractSorts(safParams)

	assert.Nil(t, sorts)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Malformed sortBy parameter. Should be <field>.<sortdirection>", err.Message())
}

func Test_extractSorts_NonexistantField_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=asdf.asc")
	safParams := url.Query()

	sorts, err := extractSorts(safParams)

	assert.Nil(t, sorts)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Unknown field asdf for sortBy", err.Message())
}

func Test_extractSorts_WrongDirection_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=id.down")
	safParams := url.Query()

	sorts, err := extractSorts(safParams)

	assert.Nil(t, sorts)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Malformed sort direction down. Should be asc or desc", err.Message())
}

func Test_extractSorts_OneSort_Returns_NoError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=id.asc")
	safParams := url.Query()

	sorts, err := extractSorts(safParams)

	assert.NotNil(t, sorts)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(sorts))
	assert.EqualValues(t, sorts[0].Field, "id")
	assert.EqualValues(t, sorts[0].Dir, "ASC")
}

func Test_extractSorts_TwoSorts_Returns_NoError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=id.asc&sortBy=rank.desc")
	safParams := url.Query()

	sorts, err := extractSorts(safParams)

	assert.NotNil(t, sorts)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(sorts))
	assert.EqualValues(t, sorts[0].Field, "id")
	assert.EqualValues(t, sorts[0].Dir, "ASC")
	assert.EqualValues(t, sorts[1].Field, "rank")
	assert.EqualValues(t, sorts[1].Dir, "DESC")
}

func Test_extractLimit_NoLimitParam_Returns_Maxlimit(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs")
	safParams := url.Query()
	maxLimit := 100

	limit, err := extractLimit(safParams, maxLimit)

	assert.NotNil(t, limit)
	assert.Nil(t, err)
	assert.EqualValues(t, maxLimit, *limit)
}

func Test_extractLimit_MalformedLimitParam_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=abc")
	safParams := url.Query()
	maxLimit := 100

	limit, err := extractLimit(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Could not convert limit abc to integer", err.Message())
}

func Test_extractLimit_LimitParamTooLow_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=-5")
	safParams := url.Query()
	maxLimit := 100

	limit, err := extractLimit(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Limit was set to -5 (too low). Must be between 1 and 100", err.Message())
}

func Test_extractLimit_LimitParamTooHigh_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=200")
	safParams := url.Query()
	maxLimit := 100

	limit, err := extractLimit(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Limit was set to 200 (too high). Must be between 1 and 100", err.Message())
}

func Test_extractLimit_LimitParamOK_Returns_Limit(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=50")
	safParams := url.Query()
	maxLimit := 100

	limit, err := extractLimit(safParams, maxLimit)

	assert.NotNil(t, limit)
	assert.Nil(t, err)
	assert.EqualValues(t, 50, *limit)
}

func Test_extractAnchor_NoAnchor_Returns_EmptyString(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs")
	safParams := url.Query()

	anchor, err := extractAnchor(safParams)

	assert.Nil(t, err)
	assert.EqualValues(t, "", anchor)
}

func Test_extractAnchor_MalformedAnchor_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?anchor=asdf")
	safParams := url.Query()

	anchor, err := extractAnchor(safParams)

	assert.NotNil(t, err)
	assert.EqualValues(t, "", anchor)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Anchor should be a ksuid", err.Message())
}

func Test_extractAnchor_WithAnchor_Returns_Anchor(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?anchor=23PF6ya4xcS9Q4of5cxFd5SlFPf")
	safParams := url.Query()

	anchor, err := extractAnchor(safParams)

	assert.Nil(t, err)
	assert.EqualValues(t, "23PF6ya4xcS9Q4of5cxFd5SlFPf", anchor)
}
