package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/service"
)

type JobUiHandler struct {
	Service service.JobService
	Cfg     *config.AppConfig
}

func NewJobUiHandler(cfg *config.AppConfig, svc service.JobService) JobUiHandler {
	return JobUiHandler{
		Cfg:     cfg,
		Service: svc,
	}
}

func (uh *JobUiHandler) JobListPage(c *gin.Context) {
	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
		Limit: 100,
	}
	jobs, _, _ := uh.Service.GetAllJobs(safReq)
	c.HTML(http.StatusOK, "joblist.page.tmpl", gin.H{
		"jobs": jobs,
	})
}

func (uh *JobUiHandler) ConfigPage(c *gin.Context) {
	configData := dto.GetConfig(uh.Cfg)
	c.HTML(http.StatusOK, "config.page.tmpl", gin.H{
		"configdata": configData,
	})
}

func (uh *JobUiHandler) AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.page.tmpl", nil)
}
