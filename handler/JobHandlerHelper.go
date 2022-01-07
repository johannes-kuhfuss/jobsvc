package handler

import (
	"fmt"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

func validateCreateJobRequest(newReq dto.CreateUpdateJobRequest) api_error.ApiErr {
	if newReq.Type == "" {
		return api_error.NewBadRequestError("Job create / update request must have a type")
	}
	if newReq.Priority != "" {
		if !domain.IsValidPriority(newReq.Priority) {
			return api_error.NewBadRequestError(fmt.Sprintf("Priority value %v does not exist", newReq.Priority))
		}
	}
	return nil
}

func validateUpdateJobRequest(newReq dto.CreateUpdateJobRequest) api_error.ApiErr {
	if newReq.Priority != "" {
		if !domain.IsValidPriority(newReq.Priority) {
			return api_error.NewBadRequestError(fmt.Sprintf("Priority value %v does not exist", newReq.Priority))
		}
	}
	return nil
}

func validateDequeueRequest(newReq dto.DequeueRequest) api_error.ApiErr {
	if newReq.Type == "" {
		return api_error.NewBadRequestError("Dequeue request must have a type")
	}
	return nil
}

func validateUpdateJobStatusRequest(newReq dto.UpdateJobStatusRequest) api_error.ApiErr {
	if newReq.Status == "" {
		return api_error.NewBadRequestError("Update status request must have a status")
	}
	if !domain.IsValidJobStatus(newReq.Status) {
		return api_error.NewBadRequestError(fmt.Sprintf("Wrong status value %v when updating job status", newReq.Status))
	}
	return nil
}

func validateUpdateJobHistoryRequest(newReq dto.UpdateJobHistoryRequest) api_error.ApiErr {
	if newReq.Message == "" {
		return api_error.NewBadRequestError("Update history request must have a message")
	}
	return nil
}
