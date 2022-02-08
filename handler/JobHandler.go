package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/segmentio/ksuid"
)

type JobHandler struct {
	Service service.JobService
	Cfg     *config.AppConfig
}

func NewJobHandler(cfg *config.AppConfig, svc service.JobService) JobHandler {
	return JobHandler{
		Cfg:     cfg,
		Service: svc,
	}
}

func (jh *JobHandler) getJobId(jobIdParam string) (string, api_error.ApiErr) {
	jobIdParam = jh.Cfg.RunTime.BmPolicy.Sanitize(jobIdParam)
	jobId, err := ksuid.Parse(jobIdParam)
	if err != nil {
		msg := "User Id should be a ksuid"
		logger.Error(msg, err)
		return "", api_error.NewBadRequestError(msg)
	}
	return jobId.String(), nil
}

func (jh *JobHandler) CreateJob(c *gin.Context) {
	var newJobReq dto.CreateUpdateJobRequest
	if err := c.ShouldBindJSON(&newJobReq); err != nil {
		msg := "Invalid JSON body in create job request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	jh.Cfg.RunTime.Sani.Sanitize(&newJobReq)
	err := validateCreateJobRequest(newJobReq)
	if err != nil {
		msg := "Could not validate input data for create job request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	result, err := jh.Service.CreateJob(newJobReq)
	if err != nil {
		logger.Error("Service error while creating job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (jh *JobHandler) GetAllJobs(c *gin.Context) {
	safParams := c.Request.URL.Query()
	safQuery, err := jh.validateSortAndFilterRequest(safParams, jh.Cfg.Misc.MaxResultLimit)
	if err != nil {
		logger.Error("Error parsing query parameters", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	jobs, totalCount, err := jh.Service.GetAllJobs(*safQuery)
	if err != nil {
		logger.Error("Service error while getting all jobs", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	countStr := fmt.Sprintf("%v", totalCount)
	c.Header("X-Total-Count", countStr)
	c.JSON(http.StatusOK, jobs)
}

func (jh *JobHandler) GetJobById(c *gin.Context) {
	jobId, err := jh.getJobId(c.Param("job_id"))
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

func (jh JobHandler) DeleteJobById(c *gin.Context) {
	jobId, err := jh.getJobId(c.Param("job_id"))
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

func (jh JobHandler) Dequeue(c *gin.Context) {
	var dqReq dto.DequeueRequest
	if err := c.ShouldBindJSON(&dqReq); err != nil {
		msg := "Invalid JSON body in dequeue request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	jh.Cfg.RunTime.Sani.Sanitize(&dqReq)
	err := validateDequeueRequest(dqReq)
	if err != nil {
		msg := "Could not validate input data for dequeue request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	result, err := jh.Service.Dequeue(dqReq)
	if err != nil {
		logger.Error("Service error while dequeuing job", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (jh JobHandler) UpdateJob(c *gin.Context) {
	jobId, err := jh.getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updJobReq dto.CreateUpdateJobRequest
	if err := c.ShouldBindJSON(&updJobReq); err != nil {
		msg := "Invalid JSON body in update job request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	jh.Cfg.RunTime.Sani.Sanitize(&updJobReq)
	err = validateUpdateJobRequest(updJobReq)
	if err != nil {
		msg := "Could not validate input data for update job request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
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

func (jh JobHandler) SetStatusById(c *gin.Context) {
	jobId, err := jh.getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updStatusReq dto.UpdateJobStatusRequest
	if err := c.ShouldBindJSON(&updStatusReq); err != nil {
		msg := "Invalid JSON body in update job status request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	jh.Cfg.RunTime.Sani.Sanitize(&updStatusReq)
	err = validateUpdateJobStatusRequest(updStatusReq)
	if err != nil {
		msg := "Could not validate input data for update job status request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	err = jh.Service.SetStatusById(jobId, updStatusReq)
	if err != nil {
		logger.Error("Service error while setting job status by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (jh JobHandler) SetHistoryById(c *gin.Context) {
	jobId, err := jh.getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err)
		return
	}
	var updHistoryReq dto.UpdateJobHistoryRequest
	if err := c.ShouldBindJSON(&updHistoryReq); err != nil {
		msg := "Invalid JSON body in update job history request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	jh.Cfg.RunTime.Sani.Sanitize(&updHistoryReq)
	err = validateUpdateJobHistoryRequest(updHistoryReq)
	if err != nil {
		msg := "Could not validate input data for update history request"
		logger.Error(msg, err)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	err = jh.Service.SetHistoryById(jobId, updHistoryReq)
	if err != nil {
		logger.Error("Service error while setting job history by id", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (jh JobHandler) DeleteAllJobs(c *gin.Context) {
	force := jh.Cfg.RunTime.BmPolicy.Sanitize(c.Query("force"))
	if force != "true" {
		msg := "Delete all jobs must be called with force=true"
		logger.Error(msg, nil)
		apiErr := api_error.NewBadRequestError(msg)
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	err := jh.Service.DeleteAllJobs()
	if err != nil {
		logger.Error("Service error while deleting all jobs", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
