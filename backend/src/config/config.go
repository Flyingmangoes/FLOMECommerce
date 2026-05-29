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
	STRIPE_CONF *StripeConfig
	ENVIRONMENT_STATUS string
}

type AppConfig struct {
	RETRY_LOGIN_COOLDOWN int
	MAX_RETRY_LOGIN int
}

type ServerConfig struct {
	HOST 			string
	PORT 			string
	ProxyHOST 		string
	ProxyPORT 		string
	FrontendHOST 	string
	FrontendPORT 	string
	JWT_SECRET 		string
	SUDO_SECRET		string
}

type DBConfig struct {
	DATABASE string
}

type RateLimitingConfig struct {
	RPM int
	BURST int
}

type StripeConfig struct {
	STRIPE_PUBLIC_KEY string
	STRIPE_SECRET_KEY string
	STRIPE_WEBHOOK_SECRET string
}

func (a *Application) Validate() error {
	if a.DB_CONF.DATABASE == ""{
		return errors.New("Database detail not found")
	}

	if a.SERV_CONF.JWT_SECRET == "" {
		return errors.New("jwt secret not found")
	}

	if a.SERV_CONF.SUDO_SECRET == "" {
		return errors.New("sudo secret not found")
	}

	if a.STRIPE_CONF.STRIPE_PUBLIC_KEY == "" {
		return errors.New("Stripe key not found")
	}

	if a.STRIPE_CONF.STRIPE_SECRET_KEY == "" {
		return errors.New("Stripe key not found")
	}

	if a.ENVIRONMENT_STATUS == "" {
		return errors.New("ENV_STATUS not found")
	}

    return nil
}

func NewConfig() *Application {
	return &Application{
		DB_CONF: &DBConfig{
			DATABASE: os.Getenv("DATABASE"),
		},
		SERV_CONF: &ServerConfig{
			HOST: os.Getenv("SERVER_HOST"),
			PORT: os.Getenv("SERVER_PORT"),
			ProxyHOST: os.Getenv("PROXY_HOST"),
			ProxyPORT: os.Getenv("PROXY_PORT"),
			FrontendHOST: os.Getenv("FRONTEND_HOST"),
			FrontendPORT: os.Getenv("FRONTEND_PORT"),
			JWT_SECRET: os.Getenv("JWT_SECRET"),
			SUDO_SECRET: os.Getenv("SUDO_SECRET"),
		},
		APP_CONF: &AppConfig{
			RETRY_LOGIN_COOLDOWN: 10,
			MAX_RETRY_LOGIN: 3,
		},
		RATE_CONF: &RateLimitingConfig{
			RPM: 30,
			BURST: 3,
		},
		STRIPE_CONF: &StripeConfig{
			STRIPE_PUBLIC_KEY: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
			STRIPE_SECRET_KEY: os.Getenv("STRIPE_SECRET_KEY"),
			STRIPE_WEBHOOK_SECRET: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		},
		ENVIRONMENT_STATUS: os.Getenv("ENVIRONMENT"),
	}
}
	