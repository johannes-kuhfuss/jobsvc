package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/config"
)

type JobUiHandler struct {
	Cfg *config.AppConfig
}

func NewJobUiHandler(cfg *config.AppConfig) JobUiHandler {
	return JobUiHandler{
		Cfg: cfg,
	}
}

func (uh *JobUiHandler) LandingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "landing.page.tmpl", nil)
}

func (uh *JobUiHandler) JobListPage(c *gin.Context) {
	c.HTML(http.StatusOK, "joblist.page.tmpl", nil)
}

func (uh *JobUiHandler) ConfigPage(c *gin.Context) {
	c.HTML(http.StatusOK, "config.page.tmpl", nil)
}

func (uh *JobUiHandler) AboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.page.tmpl", nil)
}
