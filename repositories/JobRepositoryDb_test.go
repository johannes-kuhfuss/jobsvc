package repositories

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

var (
	cfg  config.AppConfig
	jrd  JobRepositoryDb
	mock sqlmock.Sqlmock
)

func setupTest(t *testing.T) func() {
	var err error
	var db *sqlx.DB
	jrd = NewJobRepositoryDb(&cfg)
	db, mock, err = sqlmock.Newx()
	if err != nil {
		logger.Error("error creating sql mock", err)
	}
	jrd.cfg.RunTime.DbConn = db
	return func() {
		db.Close()
		jrd.cfg.RunTime.DbConn = nil
		mock = nil
	}
}

func Test_FindAll_NoStatus_Returns_DbError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v", table))).
		WillReturnError(sqlErr)

	jobs, err := jrd.FindAll("")

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error finding all jobs", err.Message())
}

func Test_FindAll_NoStatusNoResults_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v", table))).
		WillReturnRows(rows)

	jobs, err := jrd.FindAll("")

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "No jobs found", err.Message())
}

func Test_FindAll_WithStatus_Returns_Results(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	now := date.GetNowUtc()
	rows := sqlmock.NewRows([]string{"id", "correlation_id", "name", "created_at", "created_by", "modified_at", "modified_by", "status", "source", "destination", "type", "sub_type", "action", "action_details", "progress", "history", "extra_data", "priority", "rank"}).
		AddRow("23GaSImHjnOuKwdxYGP9fY8KmPC", "Corr Id 1", "Job 1", now, "me", now, "you", "running", "source 1", "destination 1", "encoding", "subtype 1", "action 1", "action details 1", 0, "2022-01-05T06:07:55Z: Job created\n", "no extra data 1", 2, 0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v WHERE status = $1", table))).
		WithArgs("running").WillReturnRows(rows)

	jobs, err := jrd.FindAll("running")

	assert.NotNil(t, jobs)
	assert.Nil(t, err)
}

func Test_FindById_DbError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v WHERE id = $1", table))).
		WithArgs("23GaSImHjnOuKwdxYGP9fY8KmPC").WillReturnError(sqlErr)

	job, err := jrd.FindById("23GaSImHjnOuKwdxYGP9fY8KmPC")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error finding job by id", err.Message())
}

func Test_FindById_NoResult_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	sqlErr := sql.ErrNoRows
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v WHERE id = $1", table))).
		WithArgs(id).WillReturnError(sqlErr)

	job, err := jrd.FindById(id)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("No job found for id %v", id), err.Message())
}

func Test_FindById_NoError_Returns_Result(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	now := date.GetNowUtc()
	row := sqlmock.NewRows([]string{"id", "correlation_id", "name", "created_at", "created_by", "modified_at", "modified_by", "status", "source", "destination", "type", "sub_type", "action", "action_details", "progress", "history", "extra_data", "priority", "rank"}).
		AddRow(id, "Corr Id 1", "Job 1", now, "me", now, "you", "running", "source 1", "destination 1", "encoding", "subtype 1", "action 1", "action details 1", 0, "2022-01-05T06:07:55Z: Job created\n", "no extra data 1", 2, 0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM %v WHERE id = $1", table))).
		WithArgs(id).WillReturnRows(row)

	job, err := jrd.FindById(id)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, id, job.Id.String())
}

func Test_Store_DbError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	job, _ := domain.NewJob("Job 1", "Encoding")
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf("INSERT INTO %v (id, correlation_id, name, created_at, created_by, modified_at, modified_by, status, source, destination, type, sub_type, action, action_details, progress, history, extra_data, priority, rank) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", table))).
		WithArgs(job.Id.String(), job.CorrelationId, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.Status, job.Source, job.Destination, job.Type, job.SubType, job.Action, job.ActionDetails, job.Progress, job.History, job.ExtraData, job.Priority, job.Rank).
		WillReturnError(sqlErr)

	err := jrd.Store(*job)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error storing new job", err.Message())
}

func Test_Store_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	job, _ := domain.NewJob("Job 1", "Encoding")
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf("INSERT INTO %v (id, correlation_id, name, created_at, created_by, modified_at, modified_by, status, source, destination, type, sub_type, action, action_details, progress, history, extra_data, priority, rank) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)", table))).
		WithArgs(job.Id.String(), job.CorrelationId, job.Name, job.CreatedAt, job.CreatedBy, job.ModifiedAt, job.ModifiedBy, job.Status, job.Source, job.Destination, job.Type, job.SubType, job.Action, job.ActionDetails, job.Progress, job.History, job.ExtraData, job.Priority, job.Rank).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := jrd.Store(*job)

	assert.Nil(t, err)
}
