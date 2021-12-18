package domain

import (
	"time"

	"github.com/segmentio/ksuid"
)

var JobStatusEnum EnumList = EnumList{[]EnumItem{{0, "created"}, {1, "queued"}, {2, "running"}, {3, "paused"}, {4, "finished"}, {5, "failed"}}}

type Job struct {
	Id            ksuid.KSUID `db:"id"`
	CorrelationId string      `db:"correlation_id"`
	Name          string      `db:"name"`
	CreatedAt     time.Time   `db:"created_at"`
	CreatedBy     string      `db:"created_by"`
	ModifiedAt    time.Time   `db:"modified_at"`
	ModifiedBy    string      `db:"modified_by"`
	Status        EnumItem    `db:"status"`
	Source        string      `db:"source"`
	Destination   string      `db:"destination"`
	Jtype         string      `db:"j_type"`
	SubType       string      `db:"sub_type"`
	Action        string      `db:"action"`
	ActionDetails string      `db:"action_details"`
	History       string      `db:"history"`
	ExtraData     string      `db:"extra_data"`
	Priority      EnumItem    `db:"priority"`
	Rank          int32       `db:"rank"`
}
