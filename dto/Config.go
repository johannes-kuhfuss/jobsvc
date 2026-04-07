package dto

import (
	"time"

	"github.com/johannes-kuhfuss/jobsvc/config"
)

type ConfigResp struct {
	ServerHost                 string
	ServerPort                 string
	ServerTlsPort              string
	ServerGracefulShutdownTime int
	ServerUseTls               bool
	ServerCertFile             string
	ServerKeyFile              string
	GinMode                    string
	DbUsername                 string
	DbHost                     string
	DbPort                     int32
	DbName                     string
	DbJobTable                 string
	MaxResultLimit             int
	StartDate                  time.Time
}

func GetConfig(cfg *config.AppConfig) ConfigResp {
	resp := ConfigResp{
		ServerHost:                 cfg.Server.Host,
		ServerPort:                 cfg.Server.Port,
		ServerTlsPort:              cfg.Server.TlsPort,
		ServerGracefulShutdownTime: cfg.Server.GracefulShutdownTime,
		ServerUseTls:               cfg.Server.UseTls,
		ServerCertFile:             cfg.Server.CertFile,
		ServerKeyFile:              cfg.Server.KeyFile,
		GinMode:                    cfg.Gin.Mode,
		DbUsername:                 cfg.Db.Username,
		DbHost:                     cfg.Db.Host,
		DbPort:                     cfg.Db.Port,
		DbName:                     cfg.Db.Name,
		DbJobTable:                 cfg.Db.JobTable,
		MaxResultLimit:             cfg.Misc.MaxResultLimit,
		StartDate:                  cfg.RunTime.StartDate,
	}
	if cfg.Server.Host == "" {
		resp.ServerHost = "localhost"
	}
	return resp
}
