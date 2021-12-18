package config

import (
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server struct {
		Host string `envconfig:"SERVER_HOST"`
		Port string `envconfig:"SERVER_PORT" default:"8080"`
	}
	Gin struct {
		Mode string `envconfig:"GIN_MODE" default:"release"`
	}
	Database struct {
		Username string `envconfig:"DB_USERNAME" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
	}
}

const (
	EnvFile = ".env"
)

var (
	Cfg Config
)

func InitConfig(file string) api_error.ApiErr {
	logger.Info("Initalizing configuration")
	loadConfig(file)
	err := envconfig.Process("", &Cfg)
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
