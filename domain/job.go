package domain

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
)

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
	Progress      int32       `db:"progress"`
	History       string      `db:"history"`
	ExtraData     string      `db:"extra_data"`
	Priority      int32       `db:"priority"`
	Rank          int32       `db:"rank"`
}

//go:generate mockgen -destination=../mocks/domain/mockJobRepository.go -package=domain github.com/johannes-kuhfuss/jobsvc/domain JobRepository
type JobRepository interface {
	Store(Job) api_error.ApiErr
	FindAll(dto.SortAndFilterRequest) (*[]Job, api_error.ApiErr)
	FindById(string) (*Job, api_error.ApiErr)
	Update(string, dto.CreateUpdateJobRequest) (*Job, api_error.ApiErr)
	DeleteById(string) api_error.ApiErr
	Dequeue(string) (*Job, api_error.ApiErr)
	SetStatusById(string, string, string) api_error.ApiErr
	SetHistoryById(string, string) api_error.ApiErr
	DeleteAllJobs() api_error.ApiErr
}

func NewJob(jobName string, jobType string) (*Job, api_error.ApiErr) {
	if strings.TrimSpace(jobType) == "" {
		return nil, api_error.NewBadRequestError("Job must have a type")
	}

	prio, _ := JobPriority.AsIndex("medium")

	newJob := Job{
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
		Progress:      0,
		History:       "",
		ExtraData:     "",
		Priority:      prio,
		Rank:          0,
	}
	newJob.AddHistory("Job created")
	return &newJob, nil
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

func (j *Job) AddHistory(msg string) {
	var sb strings.Builder
	now, _ := date.GetNowLocalString("")
	sb.WriteString(j.History)
	sb.WriteString(*now)
	sb.WriteString(": ")
	sb.WriteString(msg)
	sb.WriteString("\n")
	j.History = sb.String()
}

func (j *Job) ToJobResponseDto() dto.JobResponse {
	prio, _ := JobPriority.AsValue(j.Priority)
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
		Progress:      j.Progress,
		History:       j.History,
		ExtraData:     j.ExtraData,
		Priority:      prio,
		Rank:          j.Rank,
	}
}

func NewJobFromJobRequestDto(jobReq dto.CreateUpdateJobRequest) (*Job, api_error.ApiErr) {
	var prio int32
	newJob, err := NewJob(jobReq.Name, jobReq.Type)
	if err != nil {
		return nil, err
	}
	if jobReq.Priority != "" {
		prio, err = JobPriority.AsIndex(jobReq.Priority)
		if err != nil {
			return nil, api_error.NewBadRequestError(fmt.Sprintf("Priority value %v does not exist", jobReq.Priority))
		}
	} else {
		prio = DefaultJobPriority
	}
	newJob.CorrelationId = jobReq.CorrelationId
	newJob.Source = jobReq.Source
	newJob.Destination = jobReq.Destination
	newJob.SubType = jobReq.SubType
	newJob.Action = jobReq.Action
	newJob.ActionDetails = jobReq.ActionDetails
	newJob.ExtraData = jobReq.ExtraData
	newJob.Priority = prio
	if jobReq.Rank >= 0 {
		newJob.Rank = jobReq.Rank
	} else {
		newJob.Rank = 0
	}

	return newJob, nil
}

func GetJobDbFieldsAsStrings() []string {
	var fields []string
	val := reflect.ValueOf(Job{})
	for i := 0; i < val.Type().NumField(); i++ {
		fields = append(fields, string(val.Type().Field(i).Tag.Get("db")))
	}
	return fields
}
