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

	sort, err := extractSort(safParams)

	assert.NotNil(t, sort)
	assert.Nil(t, err)
	assert.EqualValues(t, "id", sort.Field)
	assert.EqualValues(t, "DESC", sort.Dir)
}

func Test_extractSorts_MalformedParam_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=asdf")
	safParams := url.Query()

	sort, err := extractSort(safParams)

	assert.Nil(t, sort)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Malformed sortBy parameter. Should be <field>.<sortdirection>", err.Message())
}

func Test_extractSorts_NonexistantField_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=asdf.asc")
	safParams := url.Query()

	sort, err := extractSort(safParams)

	assert.Nil(t, sort)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Unknown field asdf for sortBy", err.Message())
}

func Test_extractSorts_WrongDirection_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=id.down")
	safParams := url.Query()

	sort, err := extractSort(safParams)

	assert.Nil(t, sort)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Malformed sort direction down. Should be asc or desc", err.Message())
}

func Test_extractSorts_Returns_NoError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?sortBy=id.asc")
	safParams := url.Query()

	sort, err := extractSort(safParams)

	assert.NotNil(t, sort)
	assert.Nil(t, err)
	assert.EqualValues(t, sort.Field, "id")
	assert.EqualValues(t, sort.Dir, "ASC")
}

func Test_extractLimitAndOffset_NoLimitParam_Returns_MaxlimitZeroOffset(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.NotNil(t, limit)
	assert.NotNil(t, offset)
	assert.Nil(t, err)
	assert.EqualValues(t, maxLimit, *limit)
	assert.EqualValues(t, 0, *offset)
}

func Test_extractLimitAndOffset_MalformedLimitParam_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=abc")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.Nil(t, offset)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Could not convert limit abc to integer", err.Message())
}

func Test_extractLimitAndOffset_MalformedOffsetParam_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?offset=abc")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.Nil(t, offset)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Could not convert offset abc to integer", err.Message())
}

func Test_extractLimitAndOffset_LimitParamTooLow_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=-5")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.Nil(t, offset)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Limit was set to -5 (too low). Must be between 1 and 100", err.Message())
}

func Test_extractLimitAndOffset_LimitParamTooHigh_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=200")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.Nil(t, limit)
	assert.Nil(t, offset)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Limit was set to 200 (too high). Must be between 1 and 100", err.Message())
}

func Test_extractLimitAndOffset_ParamsOK_Returns_Params(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?limit=50&offset=10")
	safParams := url.Query()
	maxLimit := 100

	limit, offset, err := extractLimitAndOffset(safParams, maxLimit)

	assert.NotNil(t, limit)
	assert.NotNil(t, offset)
	assert.Nil(t, err)
	assert.EqualValues(t, 50, *limit)
	assert.EqualValues(t, 10, *offset)
}

func Test_extractFilters_NoFilters_Returns_EmptyResult(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.NotNil(t, filters)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(filters))
}

func Test_extractFilters_MalformedFilters_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?status=neq:1:2")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.Nil(t, filters)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Malformed filter value. Should either be single value or <operator>:<value>", err.Message())
}

func Test_extractFilters_UnknownOperator_Returns_BadRequestError(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?status=bogus:1")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.Nil(t, filters)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Unknown operator bogus for filter", err.Message())
}

func Test_extractFilters_OnlyUnknownField_Returns_EmptyResult(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?bogus=true")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.NotNil(t, filters)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(filters))
}

func Test_extractFilters_OneFieldNoOperator_Returns_ResultWithEqual(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?status=running")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.NotNil(t, filters)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(filters))
	assert.EqualValues(t, "status", filters[0].Field)
	assert.EqualValues(t, "eq", filters[0].Operator)
	assert.EqualValues(t, "running", filters[0].Value)
}

func Test_extractFilters_TwoFieldsWithOperators_Returns_Result(t *testing.T) {
	url, _ := url.Parse("http://server:8080/jobs?status=neq:running&correlation_id=ct:asdf")
	safParams := url.Query()

	filters, err := extractFilters(safParams)

	assert.NotNil(t, filters)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(filters))
}
