package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/jobsvc/dto"
	"github.com/johannes-kuhfuss/jobsvc/mocks/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

var (
	uh JobUiHandler
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

func setupUiTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockJobService(ctrl)
	uh = NewJobUiHandler(&cfg, mockService)
	router = gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("../templates/*.tmpl")
	recorder = httptest.NewRecorder()
	return func() {
		router = nil
		ctrl.Finish()
	}
}

func Test_JobListPage_Returns_Jobs(t *testing.T) {
	teardown := setupUiTest(t)
	defer teardown()
	safReq := dto.SortAndFilterRequest{
		Sorts: dto.SortBy{
			Field: "id",
			Dir:   "DESC",
		},
		Limit: 100,
	}
	dummyJobList := createDummyJobList()
	mockService.EXPECT().GetAllJobs(safReq).Return(&dummyJobList, len(dummyJobList), nil)
	router.GET("/", uh.JobListPage)
	request, _ := http.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(recorder, request)

	_, parseErr := html.Parse(io.Reader(recorder.Body))
	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.Nil(t, parseErr)
}

func Test_ConfigPage_Returns_Config(t *testing.T) {
	teardown := setupUiTest(t)
	defer teardown()
	router.GET("/config", uh.ConfigPage)
	request, _ := http.NewRequest(http.MethodGet, "/config", nil)

	router.ServeHTTP(recorder, request)

	_, parseErr := html.Parse(io.Reader(recorder.Body))
	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.Nil(t, parseErr)
}

func Test_AboutPage_Returns_About(t *testing.T) {
	teardown := setupUiTest(t)
	defer teardown()
	router.GET("/about", uh.AboutPage)
	request, _ := http.NewRequest(http.MethodGet, "/about", nil)

	router.ServeHTTP(recorder, request)

	_, parseErr := html.Parse(io.Reader(recorder.Body))
	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.Nil(t, parseErr)
}
