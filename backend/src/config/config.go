package config

import (
	"errors"
	"os"
)

type Application struct {
	DB_CONF *DBConfig
	SERV_CONF *ServerConfig
	APP_CONF *AppConfig
	RATE_CONF *RateLimitingConfig
	ENVIRONMENT_STATUS string
}

type AppConfig struct {
	CUSTOM_ALIAS_LENGTH int
	MAX_RETRY_LOGIN int
}

type ServerConfig struct {
	HOST 		string
	PORT 		string
	PrivateH string
	PrivateP string
	JWT_SECRET 	string
}

type DBConfig struct {
	DB_ADDR string
}

type RateLimitingConfig struct {
	RPM int
	BURST int
}

func (a *Application) Validate() error {
	if a.DB_CONF.DB_ADDR == ""{
		return errors.New("DB_URL is required")
	}

	if a.SERV_CONF.HOST == "" {
		return errors.New("SERV_HOST is required")
	}

	if a.SERV_CONF.PORT == "" {
		return errors.New("SERV_PORT is required")
	}

	if a.SERV_CONF.PrivateH == "" {
		return errors.New("ACCEPT_HOST is required")
	}

	if a.SERV_CONF.PrivateP == "" {
		return errors.New("ACCEPT_PORT is required")
	}

	if a.SERV_CONF.JWT_SECRET == "" {
		return errors.New("SECRET is required")
	}

	if a.ENVIRONMENT_STATUS == "" {
		return errors.New("ENV_STATUS is required")
	}

    return nil
}

func NewConfig() *Application {
	return &Application{
		DB_CONF: &DBConfig{
			DB_ADDR: os.Getenv("DB_URL"),
		},
		SERV_CONF: &ServerConfig{
			HOST: os.Getenv("SERV_HOST"),
			PORT: os.Getenv("SERV_PORT"),
			PrivateH: os.Getenv("PROXY_H"),
			PrivateP: os.Getenv("PROXY_P"),
			JWT_SECRET: os.Getenv("JWT_SECRET"),
		},
		APP_CONF: &AppConfig{
			CUSTOM_ALIAS_LENGTH: 5,
			MAX_RETRY_LOGIN: 3,
		},
		RATE_CONF: &RateLimitingConfig{
			RPM: 30,
			BURST: 3,
		},
		ENVIRONMENT_STATUS: os.Getenv("ENV_STATUS"),
	}
}
	