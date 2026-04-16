package config

import (
	"errors"
	"os"
)


type Application struct {
	DBConf *DBConfig
	ServConf *ServerConfig
	AppConf *AppConfig
	RateConf *RateLimitingConfig
	Env string
}

type AppConfig struct {
	CustomAliasLength int
}

type ServerConfig struct {
	Host string
	Port string
	JWTSecret string
}

type DBConfig struct {
	DBAddr string
}

type RateLimitingConfig struct {
	RequestPerMinute int
	Burst int
}

func (a *Application) Validate() error {
	if a.DBConf.DBAddr == ""{
		return errors.New("DB_URL is required")
	}

	if a.ServConf.Host == "" {
		return errors.New("SERV_HOST is required")
	}

	if a.ServConf.Port == "" {
		return errors.New("SERV_PORT is required")
	}

	if a.ServConf.JWTSecret == "" {
		return errors.New("SECRET is required")
	}

	if a.Env == "" {
		return errors.New("ENV_STATUS is required")
	}

    return nil
}

func NewConfig() *Application {
	return &Application{
		DBConf: &DBConfig{
			DBAddr: os.Getenv("DB_URL"),
		},
		ServConf: &ServerConfig{
			Host: os.Getenv("SERV_HOST"),
			Port: os.Getenv("SERV_PORT"),
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		AppConf: &AppConfig{
			CustomAliasLength: 6,
		},
		RateConf: &RateLimitingConfig{
			RequestPerMinute: 30,
			Burst: 3,
		},
		Env: os.Getenv("ENV_STATUS"),
	}
}
	