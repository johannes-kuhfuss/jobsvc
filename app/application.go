package app

import (
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func StartApp() {
	logger.Info("Starting application")
	mapUrls()
	logger.Info("Application ended")
}
