package domain

import (
	"net/http"
	"testing"
	"time"

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
	assert.EqualValues(t, StatusCreated, "created")
	assert.EqualValues(t, StatusQueued, "queued")
	assert.EqualValues(t, StatusRunning, "running")
	assert.EqualValues(t, StatusPaused, "paused")
	assert.EqualValues(t, StatusFinished, "finished")
	assert.EqualValues(t, StatusFailed, "failed")
	assert.EqualValues(t, PriorityRealtime, "realtime")
	assert.EqualValues(t, PriorityHigh, "high")
	assert.EqualValues(t, PriorityMedium, "medium")
	assert.EqualValues(t, PriorityLow, "low")
	assert.EqualValues(t, PriorityIdle, "idle")
}

func Test_CreateJobName_EmptyName_ReturnsGeneratedName(t *testing.T) {
	newName := createJobName("")

	assert.Contains(t, newName, "new job @")
}

func Test_CreateJobName_WithName_ReturnsName(t *testing.T) {
	newName := createJobName("jobName")

	assert.EqualValues(t, newName, "jobName")
}

func Test_NewEmptyJob_NoJobType_ReturnsBadRequestErr(t *testing.T) {
	newJob, err := NewEmptyJob("")
	assert.Nil(t, newJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Job must have a type", err.Message())
}

func Test_NewEmptyJob_WithJobType_ReturnsNewEmptyJob(t *testing.T) {
	now := date.GetNowUtc()
	jobType := "transcode"
	newJob, err := NewEmptyJob(jobType)
	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.True(t, isValidKSUID(newJob.Id))
	assert.Empty(t, newJob.CorrelationId)
	assert.Contains(t, newJob.Name, "new job @")
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
	assert.Contains(t, newJob.History.ToString(), "Job created")
	assert.Empty(t, newJob.ExtraData)
	assert.EqualValues(t, PriorityMedium, newJob.Priority)
	assert.EqualValues(t, 0, newJob.Rank)
}
