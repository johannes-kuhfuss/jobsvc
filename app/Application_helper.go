package app

import (
	"fmt"
	"time"

	"github.com/johannes-kuhfuss/services_utils/logger"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

type cleanJobs struct{}

func (c cleanJobs) Run() {
	err := jobService.CleanJobs()
	if err != nil {
		logger.Error("Error while cleaning jobs from database", nil)
	}
}
