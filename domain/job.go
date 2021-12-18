package domain

import (
	"time"

	"github.com/segmentio/ksuid"
)

type JobStatus string

const (
	Created  JobStatus = "created"
	Queued   JobStatus = "queued"
	Running  JobStatus = "running"
	Paused   JobStatus = "paused"
	Finished JobStatus = "finihsed"
	Failed   JobStatus = "failed"
)

type Job struct {
	Id         ksuid.KSUID `db:"job_id"`
	Name       string      `db:"name"`
	CreatedAt  time.Time   `db:"created_at"`
	CreatedBy  string      `db:"created_by"`
	ModifiedAt time.Time   `db:"modified_at"`
	ModifiedBy string      `db:"modified_by"`
}
