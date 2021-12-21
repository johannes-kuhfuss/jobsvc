package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/handler"
	"github.com/johannes-kuhfuss/jobsvc/repositories"
	"github.com/johannes-kuhfuss/jobsvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	cfg        config.AppConfig
	jobRepo    domain.JobRepository
	jobService service.DefaultJobService
	jobHandler handler.JobHandlers
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	initDb()
	wireApp()
	mapUrls()
	startRouter()
	cfg.RunTime.DbConn.Close()
	logger.Info("Application ended")
}

func initRouter() {
	gin.SetMode(cfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	cfg.RunTime.Router = gin.New()
	cfg.RunTime.Router.Use(gin.Logger())
	cfg.RunTime.Router.Use(gin.Recovery())
	cfg.RunTime.Router.SetTrustedProxies(nil)
}

func initDb() {
	connUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", cfg.Db.Username, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	conn, err := pgxpool.Connect(context.Background(), connUrl)
	if err != nil {
		logger.Error("Could not connect to database.", err)
		panic(err)
	}
	cfg.RunTime.DbConn = conn
}

func wireApp() {
	jobRepo = repositories.NewJobRepositoryMem()
	//jobRepo = repositories.NewJobRepositoryDb(&cfg)
	jobService = service.NewJobService(jobRepo)
	jobHandler = handler.JobHandlers{
		Service: jobService,
	}
}

func startRouter() {
	listenAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Listening on %v", listenAddr))
	if err := cfg.RunTime.Router.Run(listenAddr); err != nil {
		logger.Error("Error while starting router", err)
		panic(err)
	}
}
