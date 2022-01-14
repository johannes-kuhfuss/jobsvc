package handler

import (
	"fmt"
	"net/url"
	"strconv"
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

func (jh JobHandlers) validateSortAndFilterRequest(safParams url.Values, maxLimit int) (*dto.SortAndFilterRequest, api_error.ApiErr) {
	safReq := dto.SortAndFilterRequest{}
	sort, err := jh.extractSort(safParams)
	if err != nil {
		return nil, err
	}
	safReq.Sorts = *sort
	limit, offset, err := jh.extractLimitAndOffset(safParams, maxLimit)
	if err != nil {
		return nil, err
	}
	safReq.Limit = *limit
	safReq.Offset = *offset
	filters, err := jh.extractFilters(safParams)
	if err != nil {
		return nil, err
	}
	safReq.Filters = filters
	return &safReq, nil
}

func (jh JobHandlers) extractSort(safParams url.Values) (*dto.SortBy, api_error.ApiErr) {
	sort := dto.SortBy{}
	sortBy := safParams.Get("sortBy")
	sortBy = jh.Cfg.RunTime.BmPolicy.Sanitize(sortBy)
	if len(sortBy) == 0 {
		sort := dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		}
		return &sort, nil
	}
	sortBySplit := strings.Split(sortBy, ".")
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
	sort.Field = field
	sort.Dir = strings.ToUpper(order)
	return &sort, nil
}

func (jh JobHandlers) extractLimitAndOffset(safParams url.Values, maxLimit int) (*int, *int, api_error.ApiErr) {
	var (
		limit  int = maxLimit
		offset int = 0
		err    error
	)
	limitStr := safParams.Get("limit")
	limitStr = jh.Cfg.RunTime.BmPolicy.Sanitize(limitStr)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			msg := fmt.Sprintf("Could not convert limit %v to integer", limitStr)
			logger.Error(msg, err)
			return nil, nil, api_error.NewBadRequestError(msg)
		}
		if limit < 1 {
			msg := fmt.Sprintf("Limit was set to %v (too low). Must be between 1 and %v", limit, maxLimit)
			logger.Error(msg, nil)
			return nil, nil, api_error.NewBadRequestError(msg)
		}
		if limit > maxLimit {
			msg := fmt.Sprintf("Limit was set to %v (too high). Must be between 1 and %v", limit, maxLimit)
			logger.Error(msg, nil)
			return nil, nil, api_error.NewBadRequestError(msg)
		}
	}
	offsetStr := safParams.Get("offset")
	offsetStr = jh.Cfg.RunTime.BmPolicy.Sanitize(offsetStr)
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			msg := fmt.Sprintf("Could not convert offset %v to integer", offsetStr)
			logger.Error(msg, err)
			return nil, nil, api_error.NewBadRequestError(msg)
		}
	}
	return &limit, &offset, nil
}

func (jh JobHandlers) extractFilters(safParams url.Values) ([]dto.FilterBy, api_error.ApiErr) {
	filters := []dto.FilterBy{}
	for key, val := range safParams {
		filter := dto.FilterBy{}
		key = jh.Cfg.RunTime.BmPolicy.Sanitize(key)
		if (key != "sortBy") && (key != "limit") && (key != "offset") {
			if utils.SliceContainsString(domain.GetJobDbFieldsAsStrings(), key) {
				filter.Field = key
				for _, innerVal := range val {
					innerVal = jh.Cfg.RunTime.BmPolicy.Sanitize(innerVal)
					valSplit := strings.Split(innerVal, ":")
					if (len(valSplit) != 1) && (len(valSplit) != 2) {
						msg := "Malformed filter value. Should either be single value or <operator>:<value>"
						logger.Error(msg, nil)
						return nil, api_error.NewBadRequestError(msg)
					}
					if len(valSplit) == 1 {
						filter.Operator = "eq"
						if key == "priority" {
							prio, err := domain.JobPriority.AsIndex(valSplit[0])
							if err != nil {
								msg := fmt.Sprintf("Priority value %v does not exist", valSplit[0])
								logger.Error(msg, nil)
								return nil, api_error.NewBadRequestError(msg)
							}
							filter.Value = prio
						} else {
							filter.Value = valSplit[0]
						}

					}
					if len(valSplit) == 2 {
						if !utils.SliceContainsString(dto.Operators, valSplit[0]) {
							msg := fmt.Sprintf("Unknown operator %v for filter", valSplit[0])
							logger.Error(msg, nil)
							return nil, api_error.NewBadRequestError(msg)
						}
						filter.Operator = valSplit[0]
						if key == "priority" {
							prio, err := domain.JobPriority.AsIndex(valSplit[1])
							if err != nil {
								msg := fmt.Sprintf("Priority value %v does not exist", valSplit[1])
								logger.Error(msg, nil)
								return nil, api_error.NewBadRequestError(msg)
							}
							filter.Value = prio
						} else {
							filter.Value = valSplit[1]
						}
					}
				}
				filters = append(filters, filter)
			} else {
				logger.Info(fmt.Sprintf("Ignoring unknown filter field %v", key))
			}
		}
	}
	return filters, nil
}
