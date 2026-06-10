package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/repository"
	"backend/src/server"
	"backend/src/services"
	emailSrvc "backend/src/services/email"
	paymentSrvc "backend/src/services/payment"
	"backend/src/services/redis"
	Logger "backend/src/utils/logger"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/subosito/gotenv"
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
	Logger.LoggerInit(cfg.ENVIRONMENT_STATUS)
	defer Logger.Log.Sync()

	slog.Info("Connecting to database")
	db := database.NewDatabaseConnection(cfg.DB_CONF.DATABASE)

	slog.Info("Initializing repository")
	userStore 	 := repository.NewUserStore(db)
	storeStore 	 := repository.NewStoresStore(db)
	productStore := repository.NewProductStore(db)
	orderStore 	 := repository.NewOrderStore(db)
	tokenStore 	 := repository.NewTokenStore(db)
	cartStore 	 := repository.NewCartStore(db)

	slog.Info("Initializing Service")
	emailService := emailSrvc.NewSGMailManager(
		cfg.SENDGRID_CONF.SENDGRID_SECRET, 
		cfg.SENDGRID_CONF.VERIFICATION_SECRET,
		emailSrvc.SGTemplate{
			TEMP_EMAILCONFIRMATION: cfg.SENDGRID_CONF.TEMP_EMAILCONFIRMATION,
			TEMP_PASSRESET: cfg.SENDGRID_CONF.TEMP_PASSRESET,
		},
		userStore,
	)

	cacheService := redis.NewRedisService(
		userStore,
		productStore,
		storeStore,
		cfg,
	)

	success_url := net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/success"
    cancel_url := net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/cancel"
	paymentService := paymentSrvc.NewPaymentService(cfg.STRIPE_CONF.STRIPE_SECRET_KEY, 
		cfg.STRIPE_CONF.STRIPE_WEBHOOK_SECRET,
		success_url,
		cancel_url,
	)

	txManager := services.NewTxManager(db, productStore, orderStore, storeStore)

	slog.Info("Starting Server")
	serverManager := &server.ServerManager{
		Users: userStore,
		Stores: storeStore,
		Products: productStore,
		Orders: orderStore,
		Carts: cartStore,
		Tokens: tokenStore,
		Email: emailService,
		Payment: paymentService,
		Tx: txManager,
		Cacher: cacheService,
		JWTSecret: []byte(cfg.SERV_CONF.JWT_SECRET),
		SUDOSecret: []byte(cfg.SERV_CONF.SUDO_SECRET),
	}

	serverManager.Start(cfg)
	
	slog.Info("Shutting Down")
	os.Exit(0)
}