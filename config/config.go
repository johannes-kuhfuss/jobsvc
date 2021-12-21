package config

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Server struct {
		Host     string `envconfig:"SERVER_HOST"`
		Port     string `envconfig:"SERVER_PORT" default:"8080"`
		Shutdown bool   `ignored:"true" default:"false"`
	}
	Gin struct {
		Mode string `envconfig:"GIN_MODE" default:"release"`
	}
	Db struct {
		Username string `envconfig:"DB_USERNAME" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
		Host     string `envconfig:"DB_HOST" required:"true"`
		Port     string `envconfig:"DB_PORT" required:"true"`
		Name     string `envconfig:"DB_NAME" required:"true"`
	}
	RunTime struct {
		Router *gin.Engine
		DbConn *pgxpool.Pool
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
		return api_error.NewInternalServerError("Could not configure app, Check your environment variables", err)
	}
	logger.Info("Done initalizing configuration")
	return nil
}

func loadConfig(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		logger.Error("Could not open env file", err)
		return err
	}
	return nil
}
