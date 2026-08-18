package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/repository"
	"backend/src/server"
	"backend/src/services"
	email_service "backend/src/services/email"
	payment_service "backend/src/services/payment"
	"backend/src/services/redis"
	logger_system "backend/src/utils/LoggerSystem"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

func main() {
	slog.Info("Loading environment variable")
	if err := gotenv.Load(); err != nil {
		slog.Warn(fmt.Sprintf("detail: %v", err))
	}

	slog.Info("Loading configs")
	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		slog.Error(fmt.Sprintf("detail: %v", err))
		os.Exit(1)
	}

	slog.Info("Initializing System Logger")
	logger_system.LoggerInit(cfg.ENVIRONMENT_STATUS)
	defer logger_system.Log.Sync()

	logger_system.Log.Info("Connecting to database")
	db := database.NewDatabaseConnection(cfg.DB_CONF.DATABASE)
	
	logger_system.Log.Info("Initializing repository")
	userStore 	 := repository.NewUserStore(db)
	productStore := repository.NewProductStore(db)
	orderStore 	 := repository.NewOrderStore(db)
	tokenStore 	 := repository.NewTokenStore(db)
	cartStore 	 := repository.NewCartStore(db)

	logger_system.Log.Info("Initializing Service")

	serverUrl := net.JoinHostPort(cfg.SERV_CONF.ServerHOST, cfg.SERV_CONF.ServerPORT)

	emailService := email_service.NewSendgridManager(cfg.SENDGRID_CONF, serverUrl, userStore)
	if err := emailService.Validate(); err != nil {
		logger_system.Log.Error("Missing dependency", zap.Error(err))
		os.Exit(1)
	}

	logger_system.Log.Info("Emailing Service started")

	cacheService := cache_service.NewRedisService(cfg)
	logger_system.Log.Info("Cache Service started")

	success, cancel := payment_service.NewStripeURL(cfg)
	paymentService := payment_service.NewPaymentService(
		cfg.STRIPE_CONF, success, cancel,
	)
	logger_system.Log.Info("Payment Service started")

	txManager := services.NewTxManager(db, productStore, orderStore)

	logger_system.Log.Info("Starting Server")
	serverManager := &server.ServerManager{
		EnvironmentStatus: cfg.ENVIRONMENT_STATUS,
		Users: userStore,
		Products: productStore,	
		Orders: orderStore,
		Carts: cartStore,
		Tokens: tokenStore,
		Email: emailService,
		Payment: paymentService,
		Tx: txManager,
		Cacher: cacheService,

		ServerSecret: server.ServerSecret{
			JwtSecret: []byte(cfg.SERV_CONF.JWT_SECRET),
			SudoSecret: []byte(cfg.SERV_CONF.SUDO_SECRET),
		},
	}

	serverManager.Start(cfg)	
	logger_system.Log.Info("Shutting Down")
	os.Exit(0)
}