package repositories

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	jobRepo JobRepositoryMem
)

func setupJob() func() {
	jobRepo = NewJobRepositoryMem()
	return func() {
		jobRepo.jobList = nil
	}
}

func Test_FindAll_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	sorts := []dto.SortBy{{
		Field: "id",
		Dir:   "DESC",
	}}
	safReq := dto.SortAndFilterRequest{
		Sorts: sorts,
	}

	jList, err := jobRepo.FindAll(safReq)

	assert.Nil(t, jList)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in job list", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindAll_NoJobsAfterFilter_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()
	status := "finished"
	sorts := []dto.SortBy{{
		Field: "id",
		Dir:   "DESC",
	}}
	safReq := dto.SortAndFilterRequest{
		Sorts: sorts,
	}

	jList, err := jobRepo.FindAll(safReq)

	assert.Nil(t, jList)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("No jobs with status %v in joblist", status), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindAll_NoFilter_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()
	sorts := []dto.SortBy{{
		Field: "id",
		Dir:   "DESC",
	}}
	safReq := dto.SortAndFilterRequest{
		Sorts: sorts,
	}

	jList, err := jobRepo.FindAll(safReq)

	assert.NotNil(t, jList)
	assert.Nil(t, err)
	assert.EqualValues(t, 2, len(*jList))
}

func Test_FindAll_WithFilter_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()
	sorts := []dto.SortBy{{
		Field: "id",
		Dir:   "DESC",
	}}
	safReq := dto.SortAndFilterRequest{
		Sorts: sorts,
	}

	jList, err := jobRepo.FindAll(safReq)

	assert.NotNil(t, jList)
	assert.Nil(t, err)
	assert.NotEqual(t, jobRepo.jobList, jList)
	assert.EqualValues(t, 1, len(*jList))
}

func Test_FindById_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()

	job, err := jobRepo.FindById("")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindById_NoJobsAfterFilter_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()
	id := ksuid.New().String()

	job, err := jobRepo.FindById(id)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("No job with id %v in joblist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindById_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	id := fillJobList()

	job, err := jobRepo.FindById(id)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, id, job.Id.String())
}

func fillJobList() (id string) {
	job1, _ := domain.NewJob("job 1", "encoding")
	job2, _ := domain.NewJob("job 2", "proxy")
	job2.Status = domain.StatusRunning
	id1 := job1.Id.String()
	id2 := job2.Id.String()
	jList := make(map[string]domain.Job)
	jList[id1] = *job1
	jList[id2] = *job2
	jobRepo.mu.Lock()
	defer jobRepo.mu.Unlock()
	jobRepo.jobList = jList
	return id1
}

func Test_Store_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	job, _ := domain.NewJob("job 1", "encoding")

	err := jobRepo.Store(*job)

	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(jobRepo.jobList))
}

func Test_DeleteById_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()

	err := jobRepo.DeleteById("")

	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteById_NoJobWithId_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()
	id := ksuid.New().String()

	err := jobRepo.DeleteById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("No job with id %v in joblist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteById_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	id := fillJobList()

	deletErr := jobRepo.DeleteById(id)
	job, findErr := jobRepo.FindById(id)

	assert.Nil(t, deletErr)
	assert.NotNil(t, findErr)
	assert.Nil(t, job)
	assert.Equal(t, 1, len(jobRepo.jobList))
}

func Test_Dequeue_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()

	job, err := jobRepo.Dequeue("encoding")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_Dequeue_NoCreatedJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	createdId := fillJobList()
	jobRepo.SetStatusById(createdId, "running", "")

	job, err := jobRepo.Dequeue("")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No more jobs to dequeue", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_Dequeue_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	createdId := fillJobList()

	job, err := jobRepo.Dequeue("encoding")

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, createdId, job.Id.String())
}

func Test_SetStatusById_NoJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()

	err := jobRepo.SetStatusById("", "", "")

	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_SetStatusById_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	id := fillJobList()

	err := jobRepo.SetStatusById(id, "failed", "")

	job, _ := jobRepo.FindById(id)
	assert.Nil(t, err)
	assert.EqualValues(t, "failed", job.Status)
}

func Test_SetHistoryById_NoJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()

	err := jobRepo.SetHistoryById("", "")

	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_SetHistoryById_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	id := fillJobList()

	err := jobRepo.SetHistoryById(id, "new entry")

	job, _ := jobRepo.FindById(id)
	assert.Nil(t, err)
	assert.Contains(t, job.History, "new entry")
}

func Test_Update_NoJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	jobUpdReq := dto.CreateUpdateJobRequest{}

	job, err := jobRepo.Update("", jobUpdReq)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_Update_Returns_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	id := fillJobList()
	jobUpdReq := dto.CreateUpdateJobRequest{
		SubType: "new sub type",
		Rank:    15,
	}

	job, err := jobRepo.Update(id, jobUpdReq)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, jobUpdReq.SubType, job.SubType)
	assert.EqualValues(t, jobUpdReq.Rank, job.Rank)
}

func Test_DeleteAllJobs_NoError(t *testing.T) {
	teardown := setupJob()
	defer teardown()
	fillJobList()

	err := jobRepo.DeleteAllJobs()

	assert.Nil(t, err)
	assert.EqualValues(t, 0, len(jobRepo.jobList))
}
