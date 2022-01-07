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
