package repositories

import (
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
		findAllSql := fmt.Sprintf("SELECT * FROM %v WHERE status = ?", table)
		err = conn.Select(&jobs, findAllSql, status)
	}
	if err != nil {
		logger.Error("Database error finding all jobs", err)
		return nil, api_error.NewInternalServerError("Database error finding all jobs", err)
	}
	return &jobs, nil
}

func (jrd JobRepositoryDb) FindById(id string) (*domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) Store(job domain.Job) api_error.ApiErr {
	conn := jrd.cfg.RunTime.DbConn
	sqlInsert := "INSERT INTO joblist (id, correlation_id, name, created_at, created_by, modified_at, modified_by, status, source, destination, type, sub_type, action, action_details, history, extra_data, priority, rank) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)"
	_, err := conn.Exec(sqlInsert, job.Id.String(), job.CorrelationId, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.Status, job.Source, job.Destination, job.Type, job.SubType, job.Action, job.ActionDetails, job.History, nil, job.Priority, job.Rank)
	if err != nil {
		logger.Error("Database error storing new job", err)
		return api_error.NewInternalServerError("Database error storing new job", err)
	}
	return nil
}

func (jrd JobRepositoryDb) DeleteById(id string) api_error.ApiErr {
	panic("not implemented")
}

func (jrd JobRepositoryDb) GetNext() (*domain.Job, api_error.ApiErr) {
	panic("not implemented")
}

func (jrd JobRepositoryDb) SetStatus(id string, newStatus dto.UpdateJobStatusRequest) api_error.ApiErr {
	panic("not implemented")
}
