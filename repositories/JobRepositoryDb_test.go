package repositories

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

var (
	cfg  config.AppConfig
	jrd  JobRepositoryDb
	mock sqlmock.Sqlmock
)

type (
	AnyTime   struct{}
	AnyString struct{}
)

func (at AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (as AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

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

func Test_FindAll_NoWhere_Returns_DbError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}
	sqlErr := sql.ErrConnDone
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v ORDER BY %v %v`, table, safReq.Sorts.Field, safReq.Sorts.Dir))).
		WillReturnError(sqlErr)

	jobs, totalCount, err := jrd.FindAll(safReq)

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, totalCount)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error getting all jobs", err.Message())
}

func Test_FindAll_NoWhereNoResults_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}
	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v ORDER BY %v %v`, table, safReq.Sorts.Field, safReq.Sorts.Dir))).
		WillReturnRows(rows)

	jobs, totalCount, err := jrd.FindAll(safReq)

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, totalCount)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "No jobs found", err.Message())
}

func Test_FindAll_NoWhere_Returns_CountError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}
	now := date.GetNowUtc()
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status", "source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			"23GaSImHjnOuKwdxYGP9fY8KmPC",
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v ORDER BY %v %v`, table, safReq.Sorts.Field, safReq.Sorts.Dir))).
		WillReturnRows(rows)
	sqlErr := sql.ErrConnDone
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT count_estimate('SELECT 1 FROM %v')`, table))).WillReturnError(sqlErr)

	jobs, totalCount, err := jrd.FindAll(safReq)

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, 0, totalCount)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error getting count", err.Message())
}

func Test_FindAll_NoWhere_Returns_Results(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
	}
	now := date.GetNowUtc()
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status", "source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			"23GaSImHjnOuKwdxYGP9fY8KmPC",
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v ORDER BY %v %v`, table, safReq.Sorts.Field, safReq.Sorts.Dir))).
		WillReturnRows(rows)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT count_estimate('SELECT 1 FROM %v')`, table))).WillReturnRows(countRows)

	jobs, totalCount, err := jrd.FindAll(safReq)

	assert.NotNil(t, jobs)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, totalCount)
}

func Test_FindAll_WithWhere_Returns_Results(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{Field: "id", Dir: "DESC"},
		Filters: []dto.FilterBy{{
			Field:    "status",
			Operator: "eq",
			Value:    "running",
		}},
		Limit:  10,
		Offset: 0,
	}
	now := date.GetNowUtc()
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status", "source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			"23GaSImHjnOuKwdxYGP9fY8KmPC",
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = 'running' ORDER BY %v %v`, table, safReq.Sorts.Field, safReq.Sorts.Dir))).
		WillReturnRows(rows)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	countWhere := "status = $$running$$"
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT count_estimate('SELECT 1 FROM %v WHERE %v')`, table, countWhere))).WillReturnRows(countRows)

	jobs, totalCount, err := jrd.FindAll(safReq)

	assert.NotNil(t, jobs)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, totalCount)
}

func Test_FindById_DbError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
		WithArgs("23GaSImHjnOuKwdxYGP9fY8KmPC").WillReturnError(sqlErr)

	job, err := jrd.FindById("23GaSImHjnOuKwdxYGP9fY8KmPC")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error getting job by id", err.Message())
}

func Test_FindById_NoResult_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	sqlErr := sql.ErrNoRows
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
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
	row := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			id,
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
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
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO %v (
		id, 
		correlation_id, 
		name, 
		created_at, 
		created_by, 
		modified_at, 
		modified_by, 
		status, source, 
		destination, 
		type, 
		sub_type, 
		action, 
		action_details, 
		progress, 
		history, 
		extra_data, 
		priority, 
		rank) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`, table))).
		WithArgs(
			job.Id.String(),
			job.CorrelationId,
			job.Name,
			job.CreatedAt,
			job.CreatedBy,
			job.ModifiedAt,
			job.ModifiedBy,
			job.Status,
			job.Source,
			job.Destination,
			job.Type,
			job.SubType,
			job.Action,
			job.ActionDetails,
			job.Progress,
			job.History,
			job.ExtraData,
			job.Priority,
			job.Rank).
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
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO %v (
		id, 
		correlation_id, 
		name, 
		created_at, 
		created_by, 
		modified_at, 
		modified_by, 
		status, 
		source, 
		destination, 
		type, 
		sub_type, 
		action, 
		action_details, 
		progress, 
		history, 
		extra_data, 
		priority, 
		rank) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`, table))).
		WithArgs(
			job.Id.String(),
			job.CorrelationId,
			job.Name,
			job.CreatedAt,
			job.CreatedBy,
			job.ModifiedAt,
			job.ModifiedBy,
			job.Status,
			job.Source,
			job.Destination,
			job.Type,
			job.SubType,
			job.Action,
			job.ActionDetails,
			job.Progress,
			job.History,
			job.ExtraData,
			job.Priority,
			job.Rank).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := jrd.Store(*job)

	assert.Nil(t, err)
}

func Test_DeleteById_DbError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New()
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE id = $1`, table))).
		WithArgs(id.String()).WillReturnError(sqlErr)

	err := jrd.DeleteById(id.String())

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error deleting job by id", err.Message())
}

func Test_DeleteById_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := ksuid.New()
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE id = $1`, table))).
		WithArgs(id.String()).WillReturnResult(sqlmock.NewResult(1, 1))

	err := jrd.DeleteById(id.String())

	assert.Nil(t, err)
}

func Test_Dequeue_TransactionBeginError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	mock.ExpectBegin().WillReturnError(sqlErr)

	job, err := jrd.Dequeue("encoding")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction start error dequeuing job", err.Message())
}

func Test_Dequeue_NoJobForType_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	jobType := "encoding"
	sqlErr := sql.ErrNoRows
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = $1 AND type = $2 ORDER BY priority ASC, rank DESC limit 1`, table))).
		WithArgs(string(domain.StatusCreated), jobType).WillReturnError(sqlErr)

	job, err := jrd.Dequeue(jobType)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("No job found to dequeue for type %v", jobType), err.Message())
}

func Test_Dequeue_DbSelectError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	jobType := "encoding"
	sqlErr := sql.ErrConnDone
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = $1 AND type = $2 ORDER BY priority ASC, rank DESC limit 1`, table))).
		WithArgs(string(domain.StatusCreated), jobType).WillReturnError(sqlErr)

	job, err := jrd.Dequeue(jobType)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error dequeuing next job (select)", err.Message())
}

func Test_Dequeue_DbUpdateError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	jobType := "encoding"
	sqlErr := sql.ErrConnDone
	now := date.GetNowUtc()
	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name", "created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			id,
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = $1 AND type = $2 ORDER BY priority ASC, rank DESC limit 1`, table))).
		WithArgs(string(domain.StatusCreated), jobType).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4`, table))).
		WithArgs(AnyTime{}, "running", AnyString{}, id).WillReturnError(sqlErr)

	job, err := jrd.Dequeue(jobType)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error dequeuing next job (update)", err.Message())
}

func Test_Dequeue_TransactionCommitError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	jobType := "encoding"
	sqlErr := sql.ErrTxDone
	now := date.GetNowUtc()
	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name", "created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			id,
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = $1 AND type = $2 ORDER BY priority ASC, rank DESC limit 1`, table))).
		WithArgs(string(domain.StatusCreated), jobType).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history, progress) = ($1, $2, $3, $4) WHERE id = $5`, table))).
		WithArgs(AnyTime{}, "running", AnyString{}, 1, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sqlErr)

	job, err := jrd.Dequeue(jobType)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction end error dequeuing job", err.Message())
}

func Test_Dequeue_NoError_Returns_Job(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	jobType := "encoding"
	now := date.GetNowUtc()
	id := "23GaSImHjnOuKwdxYGP9fY8KmPC"
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name", "created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			id,
			"Corr Id 1",
			"Job 1",
			now,
			"me",
			now,
			"you",
			"running",
			"source 1",
			"destination 1",
			"encoding",
			"subtype 1",
			"action 1",
			"action details 1",
			0,
			"2022-01-05T06:07:55Z: Job created\n",
			"no extra data 1",
			20,
			0)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE status = $1 AND type = $2 ORDER BY priority ASC, rank DESC limit 1`, table))).
		WithArgs(string(domain.StatusCreated), jobType).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history, progress) = ($1, $2, $3, $4) WHERE id = $5`, table))).
		WithArgs(AnyTime{}, "running", AnyString{}, 1, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	job, err := jrd.Dequeue(jobType)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, id, job.Id.String())
}

func Test_SetStatusById_TransactionBeginError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	newStatus := "failed"
	message := "Job History Updated"
	mock.ExpectBegin().WillReturnError(sqlErr)

	err := jrd.SetStatusById(id, newStatus, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction start error setting job status by id", err.Message())
}

func Test_SetStatusById_DbSelectError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	newStatus := "failed"
	message := "Job History Updated"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT status, history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnError(sqlErr)

	err := jrd.SetStatusById(id, newStatus, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error setting job status with id (select)", err.Message())
}

func Test_SetStatusById_DbUpdateError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	newStatus := "failed"
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"status", "history"}).
		AddRow("running", "2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT status, history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4`, table))).
		WithArgs(AnyTime{}, newStatus, AnyString{}, id).WillReturnError(sqlErr)

	err := jrd.SetStatusById(id, newStatus, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error setting job status with id (update)", err.Message())
}

func Test_SetStatusById_TransactionCommitError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrTxDone
	id := ksuid.New().String()
	newStatus := "failed"
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"status", "history"}).
		AddRow("running", "2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT status, history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4`, table))).
		WithArgs(AnyTime{}, newStatus, AnyString{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sqlErr)

	err := jrd.SetStatusById(id, newStatus, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction end error setting job status by id", err.Message())
}

func Test_SetStatusById_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := ksuid.New().String()
	newStatus := "failed"
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"status", "history"}).
		AddRow("running", "2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT status, history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, status, history) = ($1, $2, $3) WHERE id = $4`, table))).
		WithArgs(AnyTime{}, newStatus, AnyString{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := jrd.SetStatusById(id, newStatus, message)

	assert.Nil(t, err)
}

func Test_Update_TransactionBeginError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "sub type",
	}
	mock.ExpectBegin().WillReturnError(sqlErr)

	job, err := jrd.Update(id, jobUpdReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction start error updating job", err.Message())
}

func Test_Update_DbSelectError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "sub type",
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnError(sqlErr)

	job, err := jrd.Update(id, jobUpdReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error updating job (select)", err.Message())
}

func Test_Update_DbUpdateError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	oldJob := domain.Job{
		Id:            ksuid.New(),
		CorrelationId: "Corr Id 1",
		Name:          "Job 1",
		CreatedAt:     time.Now().UTC(),
		CreatedBy:     "me",
		ModifiedAt:    time.Now().UTC(),
		ModifiedBy:    "you",
		Status:        "running",
		Source:        "source 1",
		Destination:   "destination 1",
		Type:          "encoding",
		SubType:       "subtype 1",
		Action:        "action 1",
		ActionDetails: "action details 1",
		Progress:      0,
		History:       "2022-01-05T06:07:55Z: Job created\n",
		ExtraData:     "no extra data 1",
		Priority:      2,
		Rank:          0,
	}

	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "subtype 2",
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			oldJob.Id.String(),
			oldJob.CorrelationId,
			oldJob.Name,
			oldJob.CreatedAt,
			oldJob.CreatedBy,
			oldJob.ModifiedAt,
			oldJob.ModifiedBy,
			oldJob.Status,
			oldJob.Source,
			oldJob.Destination,
			oldJob.Type,
			oldJob.SubType,
			oldJob.Action,
			oldJob.ActionDetails,
			oldJob.Progress,
			oldJob.History,
			oldJob.ExtraData,
			oldJob.Priority,
			oldJob.Rank)
	mergedJob := mergeJobs(&oldJob, jobUpdReq)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
		WithArgs(oldJob.Id.String()).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (
			correlation_id, 
			name, 
			modified_at, 
			modified_by, 
			source, 
			destination, 
			type, 
			sub_type, 
			action, 
			action_details, 
			history, 
			extra_data, 
			priority, 
			rank) = 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE id = $15`, table))).
		WithArgs(
			mergedJob.CorrelationId,
			mergedJob.Name,
			AnyTime{},
			mergedJob.ModifiedBy,
			mergedJob.Source,
			mergedJob.Destination,
			mergedJob.Type,
			mergedJob.SubType,
			mergedJob.Action,
			mergedJob.ActionDetails,
			mergedJob.History,
			mergedJob.ExtraData,
			mergedJob.Priority,
			mergedJob.Rank,
			oldJob.Id.String()).
		WillReturnError(sqlErr)

	job, err := jrd.Update(oldJob.Id.String(), jobUpdReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error updating job (update)", err.Message())
}

func Test_Update_TransactionCommitError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrTxDone
	oldJob := domain.Job{
		Id:            ksuid.New(),
		CorrelationId: "Corr Id 1",
		Name:          "Job 1",
		CreatedAt:     time.Now().UTC(),
		CreatedBy:     "me",
		ModifiedAt:    time.Now().UTC(),
		ModifiedBy:    "you",
		Status:        "running",
		Source:        "source 1",
		Destination:   "destination 1",
		Type:          "encoding",
		SubType:       "subtype 1",
		Action:        "action 1",
		ActionDetails: "action details 1",
		Progress:      0,
		History:       "2022-01-05T06:07:55Z: Job created\n",
		ExtraData:     "no extra data 1",
		Priority:      2,
		Rank:          0,
	}

	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "subtype 2",
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			oldJob.Id.String(),
			oldJob.CorrelationId,
			oldJob.Name,
			oldJob.CreatedAt,
			oldJob.CreatedBy,
			oldJob.ModifiedAt,
			oldJob.ModifiedBy,
			oldJob.Status,
			oldJob.Source,
			oldJob.Destination,
			oldJob.Type,
			oldJob.SubType,
			oldJob.Action,
			oldJob.ActionDetails,
			oldJob.Progress,
			oldJob.History,
			oldJob.ExtraData,
			oldJob.Priority,
			oldJob.Rank)
	mergedJob := mergeJobs(&oldJob, jobUpdReq)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
		WithArgs(oldJob.Id.String()).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (
		correlation_id, 
		name, 
		modified_at, 
		modified_by, 
		source, 
		destination, 
		type, 
		sub_type, 
		action, 
		action_details, 
		history, 
		extra_data, 
		priority, 
		rank) = 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE id = $15`, table))).
		WithArgs(
			mergedJob.CorrelationId,
			mergedJob.Name,
			AnyTime{},
			mergedJob.ModifiedBy,
			mergedJob.Source,
			mergedJob.Destination,
			mergedJob.Type,
			mergedJob.SubType,
			mergedJob.Action,
			mergedJob.ActionDetails,
			mergedJob.History,
			mergedJob.ExtraData,
			mergedJob.Priority,
			mergedJob.Rank,
			oldJob.Id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sqlErr)

	job, err := jrd.Update(oldJob.Id.String(), jobUpdReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction end error updating job", err.Message())
}

func Test_Update_NoError_Returns_Job(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	oldJob := domain.Job{
		Id:            ksuid.New(),
		CorrelationId: "Corr Id 1",
		Name:          "Job 1",
		CreatedAt:     time.Now().UTC(),
		CreatedBy:     "me",
		ModifiedAt:    time.Now().UTC(),
		ModifiedBy:    "you",
		Status:        "running",
		Source:        "source 1",
		Destination:   "destination 1",
		Type:          "encoding",
		SubType:       "subtype 1",
		Action:        "action 1",
		ActionDetails: "action details 1",
		Progress:      0,
		History:       "2022-01-05T06:07:55Z: Job created\n",
		ExtraData:     "no extra data 1",
		Priority:      2,
		Rank:          0,
	}

	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "subtype 2",
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"correlation_id",
		"name",
		"created_at",
		"created_by",
		"modified_at",
		"modified_by",
		"status",
		"source",
		"destination",
		"type",
		"sub_type",
		"action",
		"action_details",
		"progress",
		"history",
		"extra_data",
		"priority",
		"rank"}).
		AddRow(
			oldJob.Id.String(),
			oldJob.CorrelationId,
			oldJob.Name,
			oldJob.CreatedAt,
			oldJob.CreatedBy,
			oldJob.ModifiedAt,
			oldJob.ModifiedBy,
			oldJob.Status,
			oldJob.Source,
			oldJob.Destination,
			oldJob.Type,
			oldJob.SubType,
			oldJob.Action,
			oldJob.ActionDetails,
			oldJob.Progress,
			oldJob.History,
			oldJob.ExtraData,
			oldJob.Priority,
			oldJob.Rank)
	mergedJob := mergeJobs(&oldJob, jobUpdReq)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM %v WHERE id = $1`, table))).
		WithArgs(oldJob.Id.String()).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (
			correlation_id, 
			name, 
			modified_at, 
			modified_by, 
			source, 
			destination, 
			type, 
			sub_type, 
			action, 
			action_details, 
			history, 
			extra_data, 
			priority, 
			rank) = 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE id = $15`, table))).
		WithArgs(
			mergedJob.CorrelationId,
			mergedJob.Name,
			AnyTime{},
			mergedJob.ModifiedBy,
			mergedJob.Source,
			mergedJob.Destination,
			mergedJob.Type,
			mergedJob.SubType,
			mergedJob.Action,
			mergedJob.ActionDetails,
			mergedJob.History,
			mergedJob.ExtraData,
			mergedJob.Priority,
			mergedJob.Rank,
			oldJob.Id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	job, err := jrd.Update(oldJob.Id.String(), jobUpdReq)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, oldJob.Id, job.Id)
	assert.EqualValues(t, jobUpdReq.SubType, job.SubType)
}

func Test_SetHistoryById_TransactionBeginError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	message := "Job History Updated"
	mock.ExpectBegin().WillReturnError(sqlErr)

	err := jrd.SetHistoryById(id, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction start error setting job history by id", err.Message())
}

func Test_SetHistoryById_DbSelectError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	message := "Job History Updated"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnError(sqlErr)

	err := jrd.SetHistoryById(id, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error setting job history by id (select)", err.Message())
}

func Test_SetHistoryById_DbUpdateError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrConnDone
	id := ksuid.New().String()
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"history"}).
		AddRow("2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, history) = ($1, $2) WHERE id = $3`, table))).
		WithArgs(AnyTime{}, AnyString{}, id).WillReturnError(sqlErr)

	err := jrd.SetHistoryById(id, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error setting job history by id (update)", err.Message())
}

func Test_SetHistoryById_TransactionCommitError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	sqlErr := sql.ErrTxDone
	id := ksuid.New().String()
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"history"}).
		AddRow("2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, history) = ($1, $2) WHERE id = $3`, table))).
		WithArgs(AnyTime{}, AnyString{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(sqlErr)

	err := jrd.SetHistoryById(id, message)

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database transaction end error setting job history by id", err.Message())
}

func Test_SetHistoryById_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	id := ksuid.New().String()
	message := "Job History Updated"
	rows := sqlmock.NewRows([]string{"history"}).
		AddRow("2022-01-05T06:07:55Z: Job created\n")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT history FROM %v WHERE id = $1`, table))).
		WithArgs(id).WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`UPDATE %v SET (modified_at, history) = ($1, $2) WHERE id = $3`, table))).
		WithArgs(AnyTime{}, AnyString{}, id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := jrd.SetHistoryById(id, message)

	assert.Nil(t, err)
}

func Test_DeleteAllJobs_DbError_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	sqlError := sql.ErrConnDone
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v`, table))).
		WillReturnError(sqlError)

	err := jrd.DeleteAllJobs()

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error deleting all jobs", err.Message())
}

func Test_DeleteAllJobs_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v`, table))).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := jrd.DeleteAllJobs()

	assert.Nil(t, err)
}

func Test_CleanupJobs_FailedDeleteFailed_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	sqlError := sql.ErrConnDone
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'failed' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnError(sqlError)

	err := jrd.CleanupJobs()

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error deleting expired failed jobs", err.Message())
}

func Test_CleanupJobs_SucceededDeleteFailed_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	sqlError := sql.ErrConnDone
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'failed' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'finished' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnError(sqlError)

	err := jrd.CleanupJobs()

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error deleting expired succeeded jobs", err.Message())
}

func Test_CleanupJobs_InProgressFailed_Returns_InternalServerError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	sqlError := sql.ErrConnDone
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'failed' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'finished' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT count(*) FROM %v WHERE status = 'running' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnError(sqlError)

	err := jrd.CleanupJobs()

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "Database error checking for stale in-progress jobs", err.Message())
}

func Test_CleanupJobs_NoError_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'failed' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`DELETE FROM %v WHERE status = 'finished' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`SELECT count(*) FROM %v WHERE status = 'running' AND modified_at < $1`, table))).
		WithArgs(AnyTime{}).WillReturnRows(countRows)

	err := jrd.CleanupJobs()

	assert.Nil(t, err)
}
