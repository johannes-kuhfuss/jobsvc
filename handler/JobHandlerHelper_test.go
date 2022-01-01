package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/stretchr/testify/assert"
)

func Test_validateCreateUpdateJobRequest_NoType_Returns_BadRequestError(t *testing.T) {
	req := dto.CreateUpdateJobRequest{}

	err := validateCreateUpdateJobRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "job must have a type", err.Message())
}

func Test_validateCreateUpdateJobRequest_InvalidPriority_Returns_BadRequestError(t *testing.T) {
	prio := "bogus"
	req := dto.CreateUpdateJobRequest{
		Type:     "encoding",
		Priority: prio,
	}

	err := validateCreateUpdateJobRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("priority value %v does not exist", prio), err.Message())
}

func Test_validateCreateUpdateJobRequest_ValidRequest_Returns_NoError(t *testing.T) {
	req := dto.CreateUpdateJobRequest{
		Type: "encoding",
	}

	err := validateCreateUpdateJobRequest(req)

	assert.Nil(t, err)
}

func Test_validateDequeueRequest_NoType_Returns_BadRequestError(t *testing.T) {
	req := dto.DequeueRequest{}

	err := validateDequeueRequest(req)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "dequeue must have a type", err.Message())
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
	assert.EqualValues(t, "update status must have a status", err.Message())
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
	assert.EqualValues(t, "update history must have a message", err.Message())
}

func Test_validateUpdateJobHistoryRequest_Returns_NoError(t *testing.T) {
	req := dto.UpdateJobHistoryRequest{
		Message: "Hey!",
	}

	err := validateUpdateJobHistoryRequest(req)

	assert.Nil(t, err)
}