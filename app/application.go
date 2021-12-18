package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	router *gin.Engine
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile)
	if err != nil {
		panic(err)
	}
	initRouter()
	mapUrls()
	startRouter()
	logger.Info("Application ended")
}

func initRouter() {
	gin.SetMode(config.Cfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", config.Cfg.Server.Host, config.Cfg.Server.Port)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := router.Run(listenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}
