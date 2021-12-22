package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/segmentio/ksuid"
)

type JobHandlers struct {
	Service service.JobService
}

var (
	policy *bluemonday.Policy
)

func init() {
	policy = bluemonday.UGCPolicy()
}

func getJobId(jobIdParam string) (string, api_error.ApiErr) {
	jobIdParam = policy.Sanitize(jobIdParam)
	jobId, err := ksuid.Parse(jobIdParam)
	if err != nil {
		logger.Error("User Id should be a ksuid", err)
		return "", api_error.NewBadRequestError("User id should be a ksuid")
	}
	return jobId.String(), nil
}

func (jh *JobHandlers) CreateJob(c *gin.Context) {
	var newJobReq dto.CreateUpdateJobRequest
	if err := c.ShouldBindJSON(&newJobReq); err != nil {
		logger.Error("invalid JSON body in create job request", err)
		apiErr := api_error.NewBadRequestError("invalid json body for job creation")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	newJobReq.Name = policy.Sanitize(newJobReq.Name)
	result, err := jh.Service.CreateJob(newJobReq)
	if err != nil {
		logger.Error("Service error while creating job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (jh *JobHandlers) GetAllJobs(c *gin.Context) {
	status, _ := c.GetQuery("status")
	status = policy.Sanitize(status)
	jobs, err := jh.Service.GetAllJobs(status)
	if err != nil {
		logger.Error("Service error while getting all jobs", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (jh *JobHandlers) GetJobById(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	job, err := jh.Service.GetJobById(jobId)
	if err != nil {
		logger.Error("Service error while getting job by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, job)
}

func (jh JobHandlers) DeleteJobById(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	err = jh.Service.DeleteJobById(jobId)
	if err != nil {
		logger.Error("Service error while deleting job by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (jh JobHandlers) Dequeue(c *gin.Context) {
	var dqReq dto.DequeueRequest
	if err := c.ShouldBindJSON(&dqReq); err != nil {
		logger.Error("invalid JSON body in dequeue request", err)
		apiErr := api_error.NewBadRequestError("invalid json body for dequeue request")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	result, err := jh.Service.Dequeue(dqReq)
	if err != nil {
		logger.Error("Service error while dequeuing next job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (jh JobHandlers) UpdateJob(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updJobReq dto.CreateUpdateJobRequest
	if err := c.ShouldBindJSON(&updJobReq); err != nil {
		logger.Error("invalid JSON body in update job request", err)
		apiErr := api_error.NewBadRequestError("invalid json body for job update")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	result, err := jh.Service.UpdateJob(jobId, updJobReq)
	if err != nil {
		logger.Error("Service error while updating job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (jh JobHandlers) SetStatusById(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updStatusReq dto.UpdateJobStatusRequest
	if err := c.ShouldBindJSON(&updStatusReq); err != nil {
		logger.Error("invalid JSON body in update job status request", err)
		apiErr := api_error.NewBadRequestError("invalid json body for job status update")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	err = jh.Service.SetStatusById(jobId, updStatusReq)
	if err != nil {
		logger.Error("Service error while changing job status by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (jh JobHandlers) SetHistoryById(c *gin.Context) {
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updHistoryReq dto.UpdateJobHistoryRequest
	if err := c.ShouldBindJSON(&updHistoryReq); err != nil {
		logger.Error("invalid JSON body in update job history request", err)
		apiErr := api_error.NewBadRequestError("invalid json body for job history update")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	err = jh.Service.SetHistoryById(jobId, updHistoryReq)
	if err != nil {
		logger.Error("Service error while changing job history by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
