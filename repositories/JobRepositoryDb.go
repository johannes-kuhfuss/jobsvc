package repositories

import (
	"database/sql"
	"fmt"

	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
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

func (jrd JobRepositoryDb) SetStatus(id string, newStatus dto.UpdateJobStatusRequest) api_error.ApiErr {
	panic("not implemented")
}
