package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/johannes-kuhfuss/jobsvc/config"
	"github.com/johannes-kuhfuss/jobsvc/domain"
	"github.com/johannes-kuhfuss/jobsvc/handler"
	"github.com/johannes-kuhfuss/jobsvc/repositories"
	"github.com/johannes-kuhfuss/jobsvc/service"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	cfg          config.AppConfig
	jobRepo      domain.JobRepository
	jobService   service.DefaultJobService
	jobHandler   handler.JobHandler
	server       http.Server
	appEnd       chan os.Signal
	ctx          context.Context
	cancel       context.CancelFunc
	jobUiHandler handler.JobUiHandler
	bgJobs       *cron.Cron
)

func StartApp() {
	logger.Info("Starting application")
	err := config.InitConfig(config.EnvFile, &cfg)
	if err != nil {
		panic(err)
	}
	initRouter()
	initServer()
	initDb()
	initMetrics()
	wireApp()
	mapUrls()
	RegisterForOsSignals()
	createSanitizers()
	scheduleBgJobs()
	go startServer()

	<-appEnd
	cleanUp()

	if srvErr := server.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", srvErr)
	} else {
		logger.Info("Graceful shutdown finished")
	}
}

func initRouter() {
	gin.SetMode(cfg.Gin.Mode)
	gin.DefaultWriter = logger.GetLogger()
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(AddRequestId())
	router.SetTrustedProxies(nil)
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	router.LoadHTMLGlob("./templates/*.tmpl")

	cfg.RunTime.Router = router
}

func initServer() {
	var tlsConfig tls.Config

	if cfg.Server.UseTls {
		tlsConfig = tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		}
	}
	if cfg.Server.UseTls {
		cfg.RunTime.ListenAddr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.TlsPort)
	} else {
		cfg.RunTime.ListenAddr = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	}

	server = http.Server{
		Addr:              cfg.RunTime.ListenAddr,
		Handler:           cfg.RunTime.Router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
	}
	if cfg.Server.UseTls {
		server.TLSConfig = &tlsConfig
		server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}
}

func initDb() {
	logger.Info(fmt.Sprintf("Connecting to database at %v:%v", cfg.Db.Host, cfg.Db.Port))
	connUrl := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", cfg.Db.Host, cfg.Db.Port, cfg.Db.Username, cfg.Db.Password, cfg.Db.Name)
	conn, err := sqlx.Connect("postgres", connUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not connect to database at %v:%v", cfg.Db.Host, cfg.Db.Port), err)
		panic(err)
	}
	err = conn.Ping()
	if err != nil {
		logger.Error("Could not ping database", err)
		panic(err)
	}
	cfg.RunTime.DbConn = conn
	logger.Info("Successfully connected to database")
}

func initMetrics() {
	prometheusRegister()
}

func wireApp() {
	jobRepo = repositories.NewJobRepositoryDb(&cfg)
	jobService = service.NewJobService(&cfg, jobRepo)
	jobHandler = handler.NewJobHandler(&cfg, jobService)
	jobUiHandler = handler.NewJobUiHandler(&cfg, jobService)
}

func mapUrls() {
	api := cfg.RunTime.Router.Group("/jobs", validateAuth(), prometheusMetrics())
	{
		api.POST("/", jobHandler.CreateJob)
		api.GET("/", jobHandler.GetAllJobs)
		api.GET("/:job_id", jobHandler.GetJobById)
		api.DELETE("/:job_id", jobHandler.DeleteJobById)
		api.DELETE("/", jobHandler.DeleteAllJobs)
		api.PUT("/:job_id", jobHandler.UpdateJob)
		api.PUT("/:job_id/status", jobHandler.SetStatusById)
		api.PUT("/:job_id/history", jobHandler.SetHistoryById)
		api.PUT("/dequeue", jobHandler.Dequeue)

	}
	ui := cfg.RunTime.Router.Group("/")
	{
		ui.GET("/", jobUiHandler.JobListPage)
		ui.GET("/config", jobUiHandler.ConfigPage)
		ui.GET("/about", jobUiHandler.AboutPage)
	}
	cfg.RunTime.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func RegisterForOsSignals() {
	appEnd = make(chan os.Signal, 1)
	signal.Notify(appEnd, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func createSanitizers() {
	sani, err := sanitize.New()
	if err != nil {
		logger.Error("Error creating sanitizer", err)
		panic(err)
	}
	cfg.RunTime.Sani = sani
	cfg.RunTime.BmPolicy = bluemonday.UGCPolicy()
}

func scheduleBgJobs() {
	bgJobs = cron.New()
	cleanJobcycle := fmt.Sprintf("@every %dh", cfg.Cleanup.CycleHours)
	bgJobs.AddJob(cleanJobcycle, cron.NewChain(cron.DelayIfStillRunning(cron.DefaultLogger)).Then(&cleanJobs{}))
	bgJobs.Start()
}

func startServer() {
	logger.Info(fmt.Sprintf("Listening on %v", cfg.RunTime.ListenAddr))
	cfg.RunTime.StartDate = date.GetNowUtc()
	if cfg.Server.UseTls {
		if err := server.ListenAndServeTLS(cfg.Server.CertFile, cfg.Server.KeyFile); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while starting https server", err)
			panic(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error while starting http server", err)
			panic(err)
		}
	}
}

func cleanUp() {
	shutdownTime := time.Duration(cfg.Server.GracefulShutdownTime) * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), shutdownTime)
	defer func() {
		logger.Info("Cleaning up")
		bgJobs.Stop()
		cfg.RunTime.DbConn.Close()
		logger.Info("Done cleaning up")
		cancel()
	}()
}
