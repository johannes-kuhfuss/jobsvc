package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/johannes-kuhfuss/jobsvc/dto"
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

func (h *HistoryList) ToString() string {
	var history string
	for _, entry := range h.Entries {
		history = history + entry.Date.Format(date.ApiDateLayout) + ": " + entry.Message + "\n"
	}
	return history
}

func (h HistoryList) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HistoryList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &h)
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

type JobStatusUpdate struct {
	newStatus JobStatus
	errMsg    string
}

type JobRepository interface {
	Store(Job) api_error.ApiErr
	FindAll(string) (*[]Job, api_error.ApiErr)
	FindById(string) (*Job, api_error.ApiErr)
	//Search() (*[]Job, api_error.ApiErr)
	Update(string, dto.CreateUpdateJobRequest) (*Job, api_error.ApiErr)
	DeleteById(string) api_error.ApiErr
	GetNext() (*Job, api_error.ApiErr)
	SetStatus(string, dto.UpdateJobStatusRequest) api_error.ApiErr
	//AddHistory(string, HistoryItem) api_error.ApiErr
}

func NewJob(jobName string, jobType string) (*Job, api_error.ApiErr) {
	if strings.TrimSpace(jobType) == "" {
		return nil, api_error.NewBadRequestError("Job must have a type")
	}

	var history HistoryList
	now, _ := date.GetNowLocal("")
	history.Add(*now, "Job created")

	return &Job{
		Id:            ksuid.New(),
		CorrelationId: "",
		Name:          createJobName(jobName),
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

func (j *Job) ToJobResponseDto() dto.JobResponse {
	return dto.JobResponse{
		Id:            j.Id.String(),
		CorrelationId: j.CorrelationId,
		Name:          j.Name,
		CreatedAt:     j.CreatedAt,
		CreatedBy:     j.CreatedBy,
		ModifiedAt:    j.ModifiedAt,
		ModifiedBy:    j.ModifiedBy,
		Status:        string(j.Status),
		Source:        j.Source,
		Destination:   j.Destination,
		Type:          j.Type,
		SubType:       j.SubType,
		Action:        j.Action,
		ActionDetails: j.ActionDetails,
		History:       j.History.ToString(),
		ExtraData:     j.ExtraData,
		Priority:      string(j.Priority),
		Rank:          j.Rank,
	}
}

func NewJobFromJobRequestDto(jobReq dto.CreateUpdateJobRequest) (*Job, api_error.ApiErr) {
	newJob, err := NewJob(jobReq.Name, jobReq.Type)
	if err != nil {
		return nil, err
	}
	newJob.CorrelationId = jobReq.CorrelationId
	newJob.Source = jobReq.Source
	newJob.Destination = jobReq.Destination
	newJob.SubType = jobReq.SubType
	newJob.Action = jobReq.Action
	newJob.ActionDetails = jobReq.ActionDetails
	newJob.ExtraData = jobReq.ExtraData
	newJob.Priority = JobPriority(jobReq.Priority)
	if (jobReq.Rank >= 0) && (jobReq.Rank < math.MaxInt32) {
		newJob.Rank = jobReq.Rank
	} else {
		newJob.Rank = 0
	}
	return newJob, nil
}
