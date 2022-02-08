package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sanitize/sanitize"
	"github.com/jmoiron/sqlx"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/microcosm-cc/bluemonday"
)

type AppConfig struct {
	Server struct {
		Host                 string `envconfig:"SERVER_HOST" json:"serverHost"`
		Port                 string `envconfig:"SERVER_PORT" default:"8080" json:"serverPort"`
		TlsPort              string `envconfig:"SERVER_TLS_PORT" default:"8443" json:"serverTlsPort"`
		GracefulShutdownTime int    `envconfig:"GRACEFUL_SHUTDOWN_TIME" default:"10" json:"gracefulShutdownTime"`
		UseTls               bool   `envconfig:"USE_TLS" default:"false" json:"serverUseTls"`
		CertFile             string `envconfig:"CERT_FILE" default:"./cert/cert.pem" json:"serverCertFile"`
		KeyFile              string `envconfig:"KEY_FILE" default:"./cert/cert.key" json:"serverKeyFile"`
	}
	Gin struct {
		Mode string `envconfig:"GIN_MODE" default:"release" json:"ginMode"`
	}
	Db struct {
		Username string `envconfig:"DB_USERNAME" required:"true" json:"dbUserName"`
		Password string `envconfig:"DB_PASSWORD" required:"true" json:"-"`
		Host     string `envconfig:"DB_HOST" required:"true" json:"dbHost"`
		Port     int32  `envconfig:"DB_PORT" required:"true" json:"dbPort"`
		Name     string `envconfig:"DB_NAME" required:"true" json:"dbName"`
		JobTable string `envconfig:"DB_TABLE" default:"joblist" json:"dbTableName"`
	}
	Misc struct {
		MaxResultLimit int `envconfig:"MAX_RESULT_LIMIT" default:"100" json:"maxResultLimit"`
	}
	RunTime struct {
		Router     *gin.Engine         `json:"-"`
		DbConn     *sqlx.DB            `json:"-"`
		Sani       *sanitize.Sanitizer `json:"-"`
		BmPolicy   *bluemonday.Policy  `json:"-"`
		ListenAddr string              `json:"-"`
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
