package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type JobRepositoryDb struct {
	cfg *config.AppConfig
}

var (
	table string = "joblist"
)

func NewJobRepositoryDb(c *config.AppConfig) JobRepositoryDb {
	return JobRepositoryDb{c}
}

func (jrd JobRepositoryDb) FindAll(status string) (*[]domain.Job, api_error.ApiErr) {
	conn := jrd.cfg.RunTime.DbConn
	jobs := make([]domain.Job, 0)
	var err error

	if status == "" {
		findAllSql := fmt.Sprintf("SELECT * FROM %v", table)
		err = conn.Select(&jobs, findAllSql)
	} else {
		findAllSql := fmt.Sprintf("SELECT * FROM %v WHERE status = $1", table)
		err = conn.Select(&jobs, findAllSql, status)
	}
	if err != nil {
		logger.Error("Database error finding all jobs", err)
		return nil, api_error.NewInternalServerError("Database error finding all jobs", nil)
	}
	if len(jobs) == 0 {
		return nil, api_error.NewNotFoundError("No jobs found")
	}
	return &jobs, nil
}

func (jrd JobRepositoryDb) FindById(id string) (*domain.Job, api_error.ApiErr) {
	conn := jrd.cfg.RunTime.DbConn
	var job domain.Job
	findByIdSql := fmt.Sprintf("SELECT * FROM %v WHERE id = $1", table)
	err := conn.Get(&job, findByIdSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Info(fmt.Sprintf("No job found for id %v", id))
			return nil, api_error.NewNotFoundError(fmt.Sprintf("No job found for id %v", id))
		} else {
			logger.Error("Database error finding job by id", err)
			return nil, api_error.NewInternalServerError("Database error finding job by id", nil)
		}
	}
	return &job, nil
}

func (jrd JobRepositoryDb) Store(job domain.Job) api_error.ApiErr {
	conn := jrd.cfg.RunTime.DbConn
	sqlInsert := "INSERT INTO joblist (id, correlation_id, name, created_at, created_by, modified_at, modified_by, status, source, destination, type, sub_type, action, action_details, history, extra_data, priority, rank) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)"
	_, err := conn.Exec(sqlInsert, job.Id.String(), job.CorrelationId, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.Status, job.Source, job.Destination, job.Type, job.SubType, job.Action, job.ActionDetails, job.History, job.ExtraData, job.Priority, job.Rank)
	if err != nil {
		logger.Error("Database error storing new job", err)
		return api_error.NewInternalServerError("Database error storing new job", nil)
	}
	return nil
}

func (jrd JobRepositoryDb) DeleteById(id string) api_error.ApiErr {
	conn := jrd.cfg.RunTime.DbConn
	deleteByIdSql := fmt.Sprintf("DELETE FROM %v WHERE id = $1", table)
	_, err := conn.Exec(deleteByIdSql, id)
	if err != nil {
		logger.Error("Database error deleting job by id", err)
		return api_error.NewInternalServerError("Database error deleting job by id", nil)
	}
	return nil
}

func (jrd JobRepositoryDb) GetNext() (*domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) SetStatusById(id string, newStatus dto.UpdateJobStatusRequest) api_error.ApiErr {
	conn := jrd.cfg.RunTime.DbConn
	var oldJob domain.Job
	var sqlErr error
	var tx *sqlx.Tx
	tx, sqlErr = conn.Beginx()
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database transaction error updating job status with id", nil)
	}
	sqlErr = tx.Get(&oldJob, fmt.Sprintf("SELECT status, history FROM %v WHERE id = $1", table), id)
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database error updating job status with id", nil)
	}
	newHistory := oldJob.History
	if strings.TrimSpace(newStatus.Message) == "" {
		newHistory.AddNow(fmt.Sprintf("Job status changed. New status: %v", newStatus.Status))
	} else {
		newHistory.AddNow(fmt.Sprintf("Job status changed. New status: %v; %v", newStatus.Status, newStatus.Message))
	}

	sqlUpdate := "UPDATE joblist SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4"
	now := date.GetNowUtc()
	_, sqlErr = tx.Exec(sqlUpdate, now, newStatus.Status, newHistory, id)
	if sqlErr != nil {
		logger.Error("Database error updating job with id", sqlErr)
		return api_error.NewInternalServerError("Database error updating job status with id", nil)
	}
	sqlErr = tx.Commit()
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database transaction error updating job status with id", nil)
	}
	return nil
}

func (jrd JobRepositoryDb) Update(id string, jobReq dto.CreateUpdateJobRequest) (*domain.Job, api_error.ApiErr) {
	conn := jrd.cfg.RunTime.DbConn
	var oldJob domain.Job
	var sqlErr error
	var tx *sqlx.Tx

	tx, sqlErr = conn.Beginx()
	if sqlErr != nil {
		return nil, api_error.NewInternalServerError("Database transaction error updating job with id", nil)
	}
	sqlErr = tx.Get(&oldJob, fmt.Sprintf("SELECT * FROM %v WHERE id = $1", table), id)
	if sqlErr != nil {
		return nil, api_error.NewInternalServerError("Database error updating job with id", nil)
	}
	updJob, err := mergeJobs(&oldJob, jobReq)
	if err != nil {
		return nil, err
	}
	sqlUpdate := "UPDATE joblist SET (correlation_id, name, modified_at, modified_by, source, destination, type, sub_type, action, action_details, history, extra_data, priority, rank) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE id = $15"
	_, sqlErr = tx.Exec(sqlUpdate, updJob.CorrelationId, updJob.Name, updJob.ModifiedAt, updJob.ModifiedBy, updJob.Source, updJob.Destination, updJob.Type, updJob.SubType, updJob.Action, updJob.ActionDetails, updJob.History, updJob.ExtraData, updJob.Priority, updJob.Rank, updJob.Id.String())
	if sqlErr != nil {
		logger.Error("Database error updating job with id", sqlErr)
		return nil, api_error.NewInternalServerError("Database error updating job with id", nil)
	}
	sqlErr = tx.Commit()
	if sqlErr != nil {
		return nil, api_error.NewInternalServerError("Database transaction error updating job with id", nil)
	}
	return updJob, nil
}

func mergeJobs(oldJob *domain.Job, updJobReq dto.CreateUpdateJobRequest) (*domain.Job, api_error.ApiErr) {
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
	if updJobReq.ExtraData != "" {
		mergedJob.ExtraData = updJobReq.ExtraData
		changed["ExtraData"] = updJobReq.ExtraData
	} else {
		mergedJob.ExtraData = oldJob.ExtraData
	}
	if updJobReq.Priority != "" {
		mergedJob.Priority = domain.JobPriority(updJobReq.Priority)
		changed["Priority"] = string(domain.JobPriority(updJobReq.Priority))
	} else {
		mergedJob.Priority = oldJob.Priority
	}
	if updJobReq.Rank != 0 {
		mergedJob.Rank = updJobReq.Rank
		changed["Rank"] = string(updJobReq.Rank)
	} else {
		mergedJob.Rank = oldJob.Rank
	}

	newHistory := oldJob.History
	var changedStr string
	for k, v := range changed {
		changedStr = fmt.Sprintf("%v%v: %v; ", changedStr, k, v)
	}
	newHistory.AddNow(fmt.Sprintf("Job data changed. New Data: %v", changedStr))
	mergedJob.History = newHistory

	return &mergedJob, nil
}

func (jrd JobRepositoryDb) SetHistoryById(id string, newHistoryItem dto.UpdateJobHistoryRequest) api_error.ApiErr {
	conn := jrd.cfg.RunTime.DbConn
	var oldJob domain.Job
	var sqlErr error
	var tx *sqlx.Tx
	tx, sqlErr = conn.Beginx()
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database transaction error updating job history with id", nil)
	}
	sqlErr = tx.Get(&oldJob, fmt.Sprintf("SELECT history FROM %v WHERE id = $1", table), id)
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database error updating job history with id", nil)
	}
	newHistory := oldJob.History
	newHistory.AddNow(newHistoryItem.Message)
	sqlUpdate := "UPDATE joblist SET (modified_at, history) = ($1, $2) WHERE id = $3"
	now := date.GetNowUtc()
	_, sqlErr = tx.Exec(sqlUpdate, now, newHistory, id)
	if sqlErr != nil {
		logger.Error("Database error updating job history with id", sqlErr)
		return api_error.NewInternalServerError("Database error updating job history with id", nil)
	}
	sqlErr = tx.Commit()
	if sqlErr != nil {
		return api_error.NewInternalServerError("Database transaction error updating job history with id", nil)
	}
	return nil
}
