package repositories

import (
	"testing"
	"time"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func Test_mergeJobs_NoUpdates_ReturnsJob(t *testing.T) {
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
	jobUpdReq := dto.CreateUpdateJobRequest{}

	newJob := mergeJobs(&oldJob, jobUpdReq)

	assert.NotNil(t, newJob)
	assert.EqualValues(t, oldJob.Id, newJob.Id)
	assert.EqualValues(t, oldJob.CorrelationId, newJob.CorrelationId)
	assert.EqualValues(t, oldJob.Name, newJob.Name)
	assert.EqualValues(t, oldJob.CreatedAt, newJob.CreatedAt)
	assert.EqualValues(t, oldJob.CreatedBy, newJob.CreatedBy)
	assert.EqualValues(t, oldJob.Status, newJob.Status)
	assert.EqualValues(t, oldJob.Source, newJob.Source)
	assert.EqualValues(t, oldJob.Destination, newJob.Destination)
	assert.EqualValues(t, oldJob.Type, newJob.Type)
	assert.EqualValues(t, oldJob.SubType, newJob.SubType)
	assert.EqualValues(t, oldJob.Action, newJob.Action)
	assert.EqualValues(t, oldJob.ActionDetails, newJob.ActionDetails)
	assert.EqualValues(t, oldJob.Progress, newJob.Progress)
	assert.EqualValues(t, oldJob.History, newJob.History)
	assert.EqualValues(t, oldJob.ExtraData, newJob.ExtraData)
	assert.EqualValues(t, oldJob.Priority, newJob.Priority)
	assert.EqualValues(t, oldJob.Rank, newJob.Rank)
}

func Test_mergeJobs_AllUpdates_ReturnsJob(t *testing.T) {
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
		CorrelationId: "new corr id",
		Name:          "new job name",
		Source:        "new source",
		Destination:   "new destination",
		Type:          "new type",
		SubType:       "new sub type",
		Action:        "new action",
		ActionDetails: "new action details",
		ExtraData:     "new extra data",
		Priority:      "high",
		Rank:          15,
	}

	newJob := mergeJobs(&oldJob, jobUpdReq)

	assert.NotNil(t, newJob)
	assert.EqualValues(t, oldJob.Id, newJob.Id)
	assert.EqualValues(t, jobUpdReq.CorrelationId, newJob.CorrelationId)
	assert.EqualValues(t, jobUpdReq.Name, newJob.Name)
	assert.EqualValues(t, oldJob.CreatedAt, newJob.CreatedAt)
	assert.EqualValues(t, oldJob.CreatedBy, newJob.CreatedBy)
	assert.EqualValues(t, oldJob.Status, newJob.Status)
	assert.EqualValues(t, jobUpdReq.Source, newJob.Source)
	assert.EqualValues(t, jobUpdReq.Destination, newJob.Destination)
	assert.EqualValues(t, jobUpdReq.Type, newJob.Type)
	assert.EqualValues(t, jobUpdReq.SubType, newJob.SubType)
	assert.EqualValues(t, jobUpdReq.Action, newJob.Action)
	assert.EqualValues(t, jobUpdReq.ActionDetails, newJob.ActionDetails)
	assert.EqualValues(t, oldJob.Progress, newJob.Progress)
	assert.EqualValues(t, jobUpdReq.ExtraData, newJob.ExtraData)
	prio, _ := domain.JobPriority.AsIndex(jobUpdReq.Priority)
	assert.EqualValues(t, prio, newJob.Priority)
	assert.EqualValues(t, jobUpdReq.Rank, newJob.Rank)
	assert.Contains(t, newJob.History, "Job data changed. New Data:")
}
