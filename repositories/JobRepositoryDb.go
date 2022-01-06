package repositories

import (
	"database/sql"
	"fmt"

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
	sqlInsert := fmt.Sprintf("INSERT INTO %v (id, correlation_id, name, created_at, created_by, modified_at, modified_by, status, source, destination, type, sub_type, action, action_details, progress, history, extra_data, priority, rank) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", table)
	_, err := conn.Exec(sqlInsert, job.Id.String(), job.CorrelationId, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.Status, job.Source, job.Destination, job.Type, job.SubType, job.Action, job.ActionDetails, job.Progress, job.History, job.ExtraData, job.Priority, job.Rank)
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

func (jrd JobRepositoryDb) Dequeue(jobType string) (*domain.Job, api_error.ApiErr) {
	conn := jrd.cfg.RunTime.DbConn
	var nextJob domain.Job
	var sqlErr error
	var tx *sqlx.Tx

	tx, sqlErr = conn.Beginx()
	if sqlErr != nil {
		return nil, api_error.NewInternalServerError("Database transaction error dequeuing next job", nil)
	}
	sqlErr = tx.Get(&nextJob, fmt.Sprintf("SELECT * FROM %v WHERE status = $1 ORDER BY priority ASC, rank DESC limit 1", table), "created")
	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			logger.Info(fmt.Sprintf("No job found to dequeue for jobType %v", jobType))
			return nil, api_error.NewNotFoundError(fmt.Sprintf("No job found to dequeue for jobType %v", jobType))
		} else {
			return nil, api_error.NewInternalServerError("Database error dequeuing next job", nil)
		}
	}
	nextJob.AddHistory("Dequeuing job for processing")
	now := date.GetNowUtc()
	sqlUpdate := fmt.Sprintf("UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4", table)
	_, sqlErr = tx.Exec(sqlUpdate, now, "running", nextJob.History, nextJob.Id.String())
	if sqlErr != nil {
		logger.Error("Database error updating job with id", sqlErr)
		return nil, api_error.NewInternalServerError("Database error dequeuing next job", nil)
	}
	sqlErr = tx.Commit()
	if sqlErr != nil {
		return nil, api_error.NewInternalServerError("Database transaction error dequeuing next job", nil)
	}
	return &nextJob, nil
}

func (jrd JobRepositoryDb) SetStatusById(id string, newStatus string, message string) api_error.ApiErr {
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
	oldJob.AddHistory(message)
	sqlUpdate := fmt.Sprintf("UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4", table)
	now := date.GetNowUtc()
	_, sqlErr = tx.Exec(sqlUpdate, now, newStatus, oldJob.History, id)
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
	updJob := mergeJobs(&oldJob, jobReq)
	sqlUpdate := fmt.Sprintf("UPDATE %v SET (correlation_id, name, modified_at, modified_by, source, destination, type, sub_type, action, action_details, history, extra_data, priority, rank) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE id = $15", table)
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

func (jrd JobRepositoryDb) SetHistoryById(id string, message string) api_error.ApiErr {
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
	oldJob.AddHistory(message)
	sqlUpdate := fmt.Sprintf("UPDATE %v SET (modified_at, history) = ($1, $2) WHERE id = $3", table)
	now := date.GetNowUtc()
	_, sqlErr = tx.Exec(sqlUpdate, now, oldJob.History, id)
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
