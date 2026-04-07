package repositories

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/date"
)

func mergeJobs(oldJob *domain.Job, updJobReq dto.CreateUpdateJobRequest) *domain.Job {
	changed := make(map[string]string)
	mergedJob := domain.Job{}
	mergedJob.Id = oldJob.Id
	if updJobReq.CorrelationId != "" {
		mergedJob.CorrelationId = updJobReq.CorrelationId
		changed["CorrelationId"] = updJobReq.CorrelationId
	} else {
		mergedJob.CorrelationId = oldJob.CorrelationId
	}
	if updJobReq.Name != "" {
		mergedJob.Name = updJobReq.Name
		changed["Name"] = updJobReq.Name
	} else {
		mergedJob.Name = oldJob.Name
	}
	mergedJob.CreatedAt = oldJob.CreatedAt
	mergedJob.CreatedBy = oldJob.CreatedBy
	mergedJob.ModifiedAt = date.GetNowUtc()
	mergedJob.ModifiedBy = ""
	mergedJob.Status = domain.JobStatus(oldJob.Status)
	if updJobReq.Source != "" {
		mergedJob.Source = updJobReq.Source
		changed["Source"] = updJobReq.Source
	} else {
		mergedJob.Source = oldJob.Source
	}
	if updJobReq.Destination != "" {
		mergedJob.Destination = updJobReq.Destination
		changed["Destination"] = updJobReq.Destination
	} else {
		mergedJob.Destination = oldJob.Destination
	}
	if updJobReq.Type != "" {
		mergedJob.Type = updJobReq.Type
		changed["Type"] = updJobReq.Type
	} else {
		mergedJob.Type = oldJob.Type
	}
	if updJobReq.SubType != "" {
		mergedJob.SubType = updJobReq.SubType
		changed["SubType"] = updJobReq.SubType
	} else {
		mergedJob.SubType = oldJob.SubType
	}
	if updJobReq.Action != "" {
		mergedJob.Action = updJobReq.Action
		changed["Action"] = updJobReq.Action
	} else {
		mergedJob.Action = oldJob.Action
	}
	if updJobReq.ActionDetails != "" {
		mergedJob.ActionDetails = updJobReq.ActionDetails
		changed["ActionDetails"] = updJobReq.ActionDetails
	} else {
		mergedJob.ActionDetails = oldJob.ActionDetails
	}
	mergedJob.Progress = oldJob.Progress
	if updJobReq.ExtraData != "" {
		mergedJob.ExtraData = updJobReq.ExtraData
		changed["ExtraData"] = updJobReq.ExtraData
	} else {
		mergedJob.ExtraData = oldJob.ExtraData
	}
	if updJobReq.Priority != "" {
		prio, _ := domain.JobPriority.AsIndex(updJobReq.Priority)
		mergedJob.Priority = prio
		changed["Priority"] = updJobReq.Priority
	} else {
		mergedJob.Priority = oldJob.Priority
	}
	if updJobReq.Rank != 0 {
		mergedJob.Rank = updJobReq.Rank
		changed["Rank"] = string(updJobReq.Rank)
	} else {
		mergedJob.Rank = oldJob.Rank
	}

	if len(changed) > 0 {
		var changedStr string
		for k, v := range changed {
			changedStr = fmt.Sprintf("%v%v: %v; ", changedStr, k, v)
		}
		oldJob.AddHistory(fmt.Sprintf("Job data changed. New Data: %v", changedStr))
	}
	mergedJob.History = oldJob.History
	return &mergedJob
}

func constructWhereClause(safReq dto.SortAndFilterRequest) string {
	var sb strings.Builder
	for idx, where := range safReq.Filters {
		sqlFilter := dto.SqlOperatorReplacement[where.Operator]
		val := strings.Replace(sqlFilter.ValueReplace, "@@", fmt.Sprintf("%v", where.Value), -1)
		sb.WriteString(where.Field)
		sb.WriteString(" ")
		sb.WriteString(sqlFilter.SqlOperator)
		sb.WriteString(" '")
		sb.WriteString(val)
		sb.WriteString("'")
		if idx < len(safReq.Filters)-1 {
			sb.WriteString(" AND ")
		}
	}
	return sb.String()
}
