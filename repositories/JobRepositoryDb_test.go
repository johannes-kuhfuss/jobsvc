package repositories

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/services_utils/logger"
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
}
