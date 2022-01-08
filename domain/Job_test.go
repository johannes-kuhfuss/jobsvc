package domain

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func isValidKSUID(id ksuid.KSUID) bool {
	return isValidKSUIDString(id.String())
}

func isValidKSUIDString(id string) bool {
	_, parseErr := ksuid.Parse(id)
	return parseErr == nil
}

func isNowDate(t1, t2 time.Time) bool {
	t1r := t1.Round(1 * time.Minute)
	t2r := t2.Round(1 * time.Minute)
	return t1r == t2r
}

func TestConstants(t *testing.T) {
	rt, _ := JobPriority.ItemByIndex(0)
	hi, _ := JobPriority.ItemByIndex(1)
	md, _ := JobPriority.ItemByIndex(2)
	lo, _ := JobPriority.ItemByIndex(3)
	id, _ := JobPriority.ItemByIndex(4)

	assert.EqualValues(t, StatusCreated, "created")
	assert.EqualValues(t, StatusQueued, "queued")
	assert.EqualValues(t, StatusRunning, "running")
	assert.EqualValues(t, StatusPaused, "paused")
	assert.EqualValues(t, StatusFinished, "finished")
	assert.EqualValues(t, StatusFailed, "failed")
	assert.EqualValues(t, rt.Val, "realtime")
	assert.EqualValues(t, hi.Val, "high")
	assert.EqualValues(t, md.Val, "medium")
	assert.EqualValues(t, lo.Val, "low")
	assert.EqualValues(t, id.Val, "idle")
}

func Test_CreateJobName_EmptyName_ReturnsGeneratedName(t *testing.T) {
	newName := createJobName("")

	assert.Contains(t, newName, "new job @")
}

func Test_CreateJobName_WithName_ReturnsName(t *testing.T) {
	newName := createJobName("jobName")

	assert.EqualValues(t, newName, "jobName")
}

func Test_NewJob_NoJobType_ReturnsBadRequestErr(t *testing.T) {
	newJob, err := NewJob("", "")
	assert.Nil(t, newJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Job must have a type", err.Message())
}

func Test_NewJob_WithJobType_ReturnsNewJob(t *testing.T) {
	now := date.GetNowUtc()
	jobName := "New Job"
	jobType := "transcode"
	newJob, err := NewJob(jobName, jobType)
	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.True(t, isValidKSUID(newJob.Id))
	assert.Empty(t, newJob.CorrelationId)
	assert.EqualValues(t, jobName, newJob.Name)
	assert.True(t, isNowDate(newJob.CreatedAt, now))
	assert.Empty(t, newJob.CreatedBy)
	assert.True(t, isNowDate(newJob.ModifiedAt, now))
	assert.Empty(t, newJob.ModifiedBy)
	assert.EqualValues(t, StatusCreated, newJob.Status)
	assert.Empty(t, newJob.Source)
	assert.Empty(t, newJob.Destination)
	assert.EqualValues(t, jobType, newJob.Type)
	assert.Empty(t, newJob.SubType)
	assert.Empty(t, newJob.Action)
	assert.Empty(t, newJob.ActionDetails)
	assert.Contains(t, newJob.History, "Job created")
	assert.Empty(t, newJob.ExtraData)
	assert.EqualValues(t, 2, newJob.Priority)
	assert.EqualValues(t, 0, newJob.Rank)
}

func Test_ToJobResponseDto_Returns_JobResponseDto(t *testing.T) {
	jobName := "New Job"
	jobType := "transcode"
	newJob, _ := NewJob(jobName, jobType)
	fillJob(newJob)
	jobResp := newJob.ToJobResponseDto()
	prio, _ := JobPriority.AsValue(newJob.Priority)

	assert.NotNil(t, jobResp)
	assert.True(t, isValidKSUIDString(jobResp.Id))
	assert.EqualValues(t, newJob.Id.String(), jobResp.Id)
	assert.EqualValues(t, newJob.CorrelationId, jobResp.CorrelationId)
	assert.EqualValues(t, jobName, jobResp.Name)
	assert.EqualValues(t, newJob.CreatedAt, jobResp.CreatedAt)
	assert.EqualValues(t, newJob.CreatedBy, jobResp.CreatedBy)
	assert.EqualValues(t, newJob.ModifiedAt, jobResp.ModifiedAt)
	assert.EqualValues(t, newJob.ModifiedBy, jobResp.ModifiedBy)
	assert.EqualValues(t, string(newJob.Status), jobResp.Status)
	assert.EqualValues(t, newJob.Source, jobResp.Source)
	assert.EqualValues(t, newJob.Destination, jobResp.Destination)
	assert.EqualValues(t, jobType, jobResp.Type)
	assert.EqualValues(t, newJob.SubType, jobResp.SubType)
	assert.EqualValues(t, newJob.Action, jobResp.Action)
	assert.EqualValues(t, newJob.ActionDetails, jobResp.ActionDetails)
	assert.EqualValues(t, newJob.History, jobResp.History)
	assert.EqualValues(t, newJob.ExtraData, jobResp.ExtraData)
	assert.EqualValues(t, prio, jobResp.Priority)
	assert.EqualValues(t, newJob.Rank, jobResp.Rank)
}

func fillJob(job *Job) {
	job.CorrelationId = "correlation id"
	job.CreatedBy = "created by"
	job.ModifiedBy = "modified by"
	job.Source = "source"
	job.Destination = "destination"
	job.SubType = "sub type"
	job.Action = "action"
	job.ActionDetails = "action details"
	job.Rank = 25
}

func Test_NewJobFromJobRequestDto_NoType_ReturnsBadRequestError(t *testing.T) {
	newJobReq := dto.CreateUpdateJobRequest{}
	newJob, err := NewJobFromJobRequestDto(newJobReq)

	assert.Nil(t, newJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Job must have a type", err.Message())
}

func Test_NewJobFromJobRequestDto_InvalidRank_SetsRankToZero(t *testing.T) {
	newJobReq := fillJobRequest()
	newJobReq.Rank = -5
	newJob, err := NewJobFromJobRequestDto(newJobReq)

	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, newJob.Rank)
}

func Test_NewJobFromJobRequestDto_InvalidPriority_Returns_BadRequestError(t *testing.T) {
	newJobReq := fillJobRequest()
	newJobReq.Priority = "bogus"
	newJob, err := NewJobFromJobRequestDto(newJobReq)

	assert.Nil(t, newJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, fmt.Sprintf("Priority value %v does not exist", newJobReq.Priority), err.Message())
}

func Test_NewJobFromJobRequestDto_EmptyPriority_Returns_DefaultPriorty(t *testing.T) {
	newJobReq := fillJobRequest()
	newJobReq.Priority = ""
	newJob, err := NewJobFromJobRequestDto(newJobReq)

	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.EqualValues(t, DefaultJobPriority, newJob.Priority)
}

func fillJobRequest() dto.CreateUpdateJobRequest {
	return dto.CreateUpdateJobRequest{
		CorrelationId: "corr id",
		Name:          "my new job",
		Source:        "source",
		Destination:   "destination",
		Type:          "testjob",
		SubType:       "subtype",
		Action:        "action",
		ActionDetails: "action details",
		ExtraData:     "extra data",
		Priority:      "High",
		Rank:          25,
	}
}

func Test_NewJobFromJobRequestDto_ValidValues(t *testing.T) {
	newJobReq := fillJobRequest()
	newJob, err := NewJobFromJobRequestDto(newJobReq)
	prio, _ := JobPriority.AsIndex(newJobReq.Priority)

	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.EqualValues(t, newJobReq.CorrelationId, newJob.CorrelationId)
	assert.EqualValues(t, newJobReq.Name, newJob.Name)
	assert.EqualValues(t, newJobReq.Source, newJob.Source)
	assert.EqualValues(t, newJobReq.Destination, newJob.Destination)
	assert.EqualValues(t, newJobReq.Type, newJob.Type)
	assert.EqualValues(t, newJobReq.SubType, newJob.SubType)
	assert.EqualValues(t, newJobReq.Action, newJob.Action)
	assert.EqualValues(t, newJobReq.ActionDetails, newJob.ActionDetails)
	assert.EqualValues(t, newJobReq.ExtraData, newJob.ExtraData)
	assert.EqualValues(t, prio, newJob.Priority)
	assert.EqualValues(t, newJobReq.Rank, newJob.Rank)
}

func Test_GetJobDbFieldsAsStrings(t *testing.T) {
	expectedFields := []string{"id", "correlation_id", "name", "created_at", "created_by", "modified_at", "modified_by", "status",
		"source", "destination", "type", "sub_type", "action", "action_details", "progress", "history", "extra_data", "priority", "rank"}

	jobFields := GetJobDbFieldsAsStrings()

	assert.Equal(t, expectedFields, jobFields)
}
