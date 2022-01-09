package handler

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/utils"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
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

func validateSortAndFilterRequest(safParams url.Values) (*dto.SortAndFilterRequest, api_error.ApiErr) {
	safReq := dto.SortAndFilterRequest{}
	sorts, err := extractSorts(safParams)
	if err != nil {
		return nil, err
	}
	safReq.Sorts = sorts
	return &safReq, nil
}

func extractSorts(safParams url.Values) ([]dto.SortBy, api_error.ApiErr) {
	sorts := []dto.SortBy{}
	sortBy := safParams["sortBy"]
	if len(sortBy) == 0 {
		sort := dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		}
		sorts = append(sorts, sort)
		return sorts, nil
	}
	for _, val := range sortBy {
		sortBySplit := strings.Split(val, ".")
		if len(sortBySplit) != 2 {
			msg := "Malformed sortBy parameter. Should be <field>.<sortdirection>"
			logger.Error(msg, nil)
			return nil, api_error.NewBadRequestError(msg)
		}
		field := sortBySplit[0]
		order := strings.ToLower(sortBySplit[1])
		if !utils.SliceContainsString(domain.GetJobDbFieldsAsStrings(), field) {
			msg := fmt.Sprintf("Unknown field %v for sortBy", field)
			logger.Error(msg, nil)
			return nil, api_error.NewBadRequestError(msg)
		}
		if order != "desc" && order != "asc" {
			msg := fmt.Sprintf("Malformed sort direction %v. Should be asc or desc", order)
			logger.Error(msg, nil)
			return nil, api_error.NewBadRequestError(msg)
		}
		sorts = append(sorts, dto.SortBy{
			Field: field,
			Dir:   strings.ToUpper(order),
		})
	}
	return sorts, nil
}
