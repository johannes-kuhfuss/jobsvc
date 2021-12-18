package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
)

type JobStatus string
type JobPriority string

const (
	StatusCreated  JobStatus = "created"
	StatusQueued   JobStatus = "queued"
	StatusRunning  JobStatus = "running"
	StatusPaused   JobStatus = "paused"
	StatusFinished JobStatus = "finished"
	StatusFailed   JobStatus = "failed"
)

const (
	PriorityRealtime JobPriority = "realtime"
	PriorityHigh     JobPriority = "high"
	PriorityMedium   JobPriority = "medium"
	PriorityLow      JobPriority = "low"
	PriorityIdle     JobPriority = "idle"
)

type HistoryItem struct {
	Date    time.Time
	Message string
}

type HistoryList struct {
	Entries []HistoryItem
}

func (h *HistoryList) Add(date time.Time, msg string) {
	var newEntry HistoryItem
	newEntry.Date = date
	newEntry.Message = msg
	h.Entries = append(h.Entries, newEntry)
}

type Job struct {
	Id            ksuid.KSUID `db:"id"`
	CorrelationId string      `db:"correlation_id"`
	Name          string      `db:"name"`
	CreatedAt     time.Time   `db:"created_at"`
	CreatedBy     string      `db:"created_by"`
	ModifiedAt    time.Time   `db:"modified_at"`
	ModifiedBy    string      `db:"modified_by"`
	Status        JobStatus   `db:"status"`
	Source        string      `db:"source"`
	Destination   string      `db:"destination"`
	Type          string      `db:"type"`
	SubType       string      `db:"sub_type"`
	Action        string      `db:"action"`
	ActionDetails string      `db:"action_details"`
	History       HistoryList `db:"history"`
	ExtraData     string      `db:"extra_data"`
	Priority      JobPriority `db:"priority"`
	Rank          int32       `db:"rank"`
}

func NewEmptyJob(jobType string) (*Job, api_error.ApiErr) {
	if strings.TrimSpace(jobType) == "" {
		return nil, api_error.NewBadRequestError("Job must have a type")
	}

	var history HistoryList
	now, _ := date.GetNowLocal("")
	history.Add(*now, "Job created")

	return &Job{
		Id:            ksuid.New(),
		CorrelationId: "",
		Name:          createJobName(""),
		CreatedAt:     date.GetNowUtc(),
		CreatedBy:     "",
		ModifiedAt:    date.GetNowUtc(),
		ModifiedBy:    "",
		Status:        StatusCreated,
		Source:        "",
		Destination:   "",
		Type:          jobType,
		SubType:       "",
		Action:        "",
		ActionDetails: "",
		History:       history,
		ExtraData:     "",
		Priority:      PriorityMedium,
		Rank:          0,
	}, nil
}

func createJobName(name string) string {
	var jobName string
	if strings.TrimSpace(name) == "" {
		newDate, _ := date.GetNowLocalString("")
		jobName = fmt.Sprintf("new job @ %s", *newDate)
	} else {
		jobName = name
	}
	return jobName
}
