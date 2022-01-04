package repositories

import (
	"database/sql"
	"net/http"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/jobsvc/config"
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
	jrd.cfg = &cfg
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM joblist")).WillReturnError(sqlErr)

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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM joblist")).WillReturnRows(rows)

	jobs, err := jrd.FindAll("")

	assert.Nil(t, jobs)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "No jobs found", err.Message())
}

/*
func Test_FindAll_WithStatusNoResults_Returns_Results(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// define result set
	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM joblist")).WillReturnRows(rows)

	jobs, err := jrd.FindAll("")

	assert.NotNil(t, jobs)
	assert.Nil(t, err)
}
*/
