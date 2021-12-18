package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var appCfg config.AppConfig

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &appCfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	mapUrls()
	startRouter()
	logger.Info("Application ended")
}

func initRouter() {
	gin.SetMode(appCfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	appCfg.RunTime.Router = gin.New()
	appCfg.RunTime.Router.Use(gin.Logger())
	appCfg.RunTime.Router.Use(gin.Recovery())
	appCfg.RunTime.Router.SetTrustedProxies(nil)
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", appCfg.Server.Host, appCfg.Server.Port)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := appCfg.RunTime.Router.Run(listenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}
