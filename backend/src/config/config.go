package config

import (
	"errors"
	"os"
)

type ConfigManager struct {
	DB_CONF 			*DBConfig
	SERV_CONF 			*ServerConfig
	APP_CONF 			*AppConfig
	RATE_CONF 			*RateLimitingConfig
	STRIPE_CONF 		*StripeConfig
	SENDGRID_CONF 		*SendgridConfig
	REDIS_CONF 			*RedisConfig
	ENVIRONMENT_STATUS 	string
}

type AppConfig struct {
	RETRY_LOGIN_COOLDOWN 	int
	MAX_RETRY_LOGIN 		int
}

type ServerConfig struct {
	HOST 				string
	PORT 				string

	ProxyHOST 			string
	ProxyPORT 			string

	FrontendHOST 		string
	FrontendPORT 		string

	JWT_SECRET 			string
	SUDO_SECRET			string
}

type DBConfig struct {
	DATABASE 	string
}

type RateLimitingConfig struct {
	RPM 	int
	BURST 	int
}

type StripeConfig struct {
	STRIPE_PUBLIC_KEY 		string
	STRIPE_SECRET_KEY 		string
	STRIPE_WEBHOOK_SECRET 	string
}

type SendgridConfig struct {
	VERIFICATION_SECRET 		string
	SENDGRID_SECRET 			string
	TEST_EMAIL					string
	DOMAIN_EMAIL				string

	TEMPLATE_USER_VERIFICATION 	string
	TEMPLATE_PASS_RESET 		string
	TEMPLATE_ORDER_CONFIRMATION	string
}

type RedisConfig struct {
	REDIS_PORT 	string
	REDIS_HOST 	string
	CACHE_TTL 	int
}

func (a *ConfigManager) Validate() error {
	if a.DB_CONF.DATABASE == ""{
		return errors.New("Database detail not found")
	}

	if a.SENDGRID_CONF.SENDGRID_SECRET == "" {
		return errors.New("sendgrid secret not found")
	}

	if a.SERV_CONF.JWT_SECRET == "" {
		return errors.New("jwt secret not found")
	}

	if a.SERV_CONF.SUDO_SECRET == "" {
		return errors.New("sudo secret not found")
	}

	if a.SENDGRID_CONF.VERIFICATION_SECRET == "" {
		return errors.New("verification secret not found")
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

func NewConfig() *ConfigManager {
	return &ConfigManager{
		DB_CONF: &DBConfig{
			DATABASE: os.Getenv("POSTGRES_URL"),
		},
		SERV_CONF: &ServerConfig{
			HOST: 							os.Getenv("SERVER_HOST"),
			PORT: 							os.Getenv("SERVER_PORT"),
			ProxyHOST: 						os.Getenv("PROXY_HOST"),
			ProxyPORT: 						os.Getenv("PROXY_PORT"),
			FrontendHOST: 					os.Getenv("FRONTEND_HOST"),
			FrontendPORT: 					os.Getenv("FRONTEND_PORT"),
			JWT_SECRET: 					os.Getenv("JWT_SECRET"),
			SUDO_SECRET: 					os.Getenv("SUDO_SECRET"),
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
			STRIPE_PUBLIC_KEY: 				os.Getenv("STRIPE_PUBLISHABLE_KEY"),
			STRIPE_SECRET_KEY: 				os.Getenv("STRIPE_SECRET_KEY"),
			STRIPE_WEBHOOK_SECRET: 			os.Getenv("STRIPE_WEBHOOK_SECRET"),
		},
		SENDGRID_CONF: &SendgridConfig{
			SENDGRID_SECRET: 				os.Getenv("SENDGRID_SECRET"),
			VERIFICATION_SECRET: 			os.Getenv("VERIFICATION_SECRET"),
			TEST_EMAIL: 					os.Getenv("SENDGRID_TEST_EMAIL"),
			DOMAIN_EMAIL: 					os.Getenv("SENDGRID_DOMAIN_EMAIL"),
			TEMPLATE_USER_VERIFICATION: 	os.Getenv("USER_VERIFICATION_TEMPLATE"),
			TEMPLATE_ORDER_CONFIRMATION: 	os.Getenv("ORDER_CONFIRMATION_TEMPLATE"),
			TEMPLATE_PASS_RESET: 			os.Getenv("PASS_RESET_TEMPLATE"),
		},
		ENVIRONMENT_STATUS: os.Getenv("ENVIRONMENT"),
		REDIS_CONF: &RedisConfig{
			REDIS_PORT: 					os.Getenv("REDIS_PORT"),
			REDIS_HOST: 					os.Getenv("REDIS_HOST"),
			CACHE_TTL: 5,
		},
	}
}
	