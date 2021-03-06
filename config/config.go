package config

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/microcosm-cc/bluemonday"
)

type AppConfig struct {
	Server struct {
		Host                 string `envconfig:"SERVER_HOST"`
		Port                 string `envconfig:"SERVER_PORT" default:"8080"`
		TlsPort              string `envconfig:"SERVER_TLS_PORT" default:"8443"`
		GracefulShutdownTime int    `envconfig:"GRACEFUL_SHUTDOWN_TIME" default:"10"`
		UseTls               bool   `envconfig:"USE_TLS" default:"false"`
		CertFile             string `envconfig:"CERT_FILE" default:"./cert/cert.pem"`
		KeyFile              string `envconfig:"KEY_FILE" default:"./cert/cert.key"`
	}
	Gin struct {
		Mode string `envconfig:"GIN_MODE" default:"release"`
	}
	Db struct {
		Username string `envconfig:"DB_USERNAME" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
		Host     string `envconfig:"DB_HOST" required:"true"`
		Port     int32  `envconfig:"DB_PORT" required:"true"`
		Name     string `envconfig:"DB_NAME" required:"true"`
		JobTable string `envconfig:"DB_TABLE" default:"joblist"`
	}
	Misc struct {
		MaxResultLimit int      `envconfig:"MAX_RESULT_LIMIT" default:"100"`
		ApiKeys        []string `envconfig:"API_KEYS"`
	}
	Cleanup struct {
		CycleHours             int `envconfig:"CLEANUP_CYCLE_HOURS" default:"1"`
		FailedRetentionDays    int `envconfig:"CLEANUP_FAILED_RETEN_DAYS" default:"2"`
		SuccessRetentionDays   int `envconfig:"CLEANUP_SUCCESS_RETEN_DAYS" default:"1"`
		InProgressWarningHours int `envconfig:"IN_PROGRESS_WARNING_HOURS" default:"6"`
	}
	RunTime struct {
		Router     *gin.Engine
		DbConn     *sqlx.DB
		Sani       *sanitize.Sanitizer
		BmPolicy   *bluemonday.Policy
		ListenAddr string
		StartDate  time.Time
	}
}

const (
	EnvFile = ".env"
)

func InitConfig(file string, config *AppConfig) api_error.ApiErr {
	logger.Info("Initalizing configuration")
	loadConfig(file)
	err := envconfig.Process("", config)
	if err != nil {
		return api_error.NewInternalServerError("Could not initalize configuration. Check your environment variables", err)
	}
	if len(config.Misc.ApiKeys) == 0 {
		id, _ := uuid.NewV4()
		config.Misc.ApiKeys = append(config.Misc.ApiKeys, id.String())
	}
	logger.Info("Done initalizing configuration")
	return nil
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Info("Could not open env file. Using Environment variable and defaults")
		return err
	}
	return nil
}
