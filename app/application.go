package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/handler"
	"github.com/johannes-kuhfuss/jobsvc/repositories"
	"github.com/johannes-kuhfuss/jobsvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/microcosm-cc/bluemonday"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	if cfg.BackEnd == "db" {
		initDb()
	}
	wireApp()
	mapUrls()
	createSanitizer()
	createBmPolicy()
	startRouter()
	cfg.RunTime.DbConn.Close()
	logger.Info("Application ended")
}

func createSanitizer() {
	sani, err := sanitize.New()
	if err != nil {
		logger.Error("error creating sanitizer", err)
		panic(err)
	}
	cfg.RunTime.Sani = sani
}

func createBmPolicy() {
	cfg.RunTime.BmPolicy = bluemonday.UGCPolicy()
}

func initRouter() {
	gin.SetMode(cfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies(nil)
	cfg.RunTime.Router = router
}

func initDb() {
	logger.Info(fmt.Sprintf("Connecting to database at %v:%v", cfg.Db.Host, cfg.Db.Port))
	connUrl := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.Db.Host, cfg.Db.Port, cfg.Db.Username, cfg.Db.Password, cfg.Db.Name)
	conn, err := sqlx.Connect("postgres", connUrl)
	if err != nil {
		logger.Error("Could not connect to database", err)
		panic(err)
	}
	cfg.RunTime.DbConn = conn
	logger.Info("Successfully connected to database")
}

func wireApp() {
	if cfg.BackEnd == "db" {
		jobRepo = repositories.NewJobRepositoryDb(&cfg)
	} else {
		jobRepo = repositories.NewJobRepositoryMem()
	}
	jobService = service.NewJobService(jobRepo)
	jobHandler = handler.JobHandlers{
		Service: jobService,
		Cfg:     &cfg,
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
