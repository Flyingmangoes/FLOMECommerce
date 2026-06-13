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
	Logger.LoggerInit(cfg.ENVIRONMENT_STATUS)
	defer Logger.Log.Sync()

	Logger.Log.Info("Connecting to database")
	db := database.NewDatabaseConnection(cfg.DB_CONF.DATABASE)
	
	Logger.Log.Info("Initializing repository")
	userStore 	 := repository.NewUserStore(db)
	storeStore 	 := repository.NewStoresStore(db)
	productStore := repository.NewProductStore(db)
	orderStore 	 := repository.NewOrderStore(db)
	tokenStore 	 := repository.NewTokenStore(db)
	cartStore 	 := repository.NewCartStore(db)

	Logger.Log.Info("Initializing Service")
	emailService := emailSrvc.NewSGMailManager(
		cfg.SENDGRID_CONF.SENDGRID_SECRET, 
		cfg.SENDGRID_CONF.VERIFICATION_SECRET,
		cfg.SENDGRID_CONF.TEST_EMAIL,
		cfg.SENDGRID_CONF.DOMAIN_EMAIL,
		emailSrvc.SGTemplate{
			T_UserVerification: cfg.SENDGRID_CONF.TEMPLATE_USER_VERIFICATION,
			T_PassReset: cfg.SENDGRID_CONF.TEMPLATE_PASS_RESET,
			T_OrderConfirmation: cfg.SENDGRID_CONF.TEMPLATE_ORDER_CONFIRMATION,
		},
		userStore,
	)

	if err := emailService.Validate(); err != nil {
		Logger.Log.Error("Missing dependency", zap.Error(err))
		os.Exit(1)
	}
	
	Logger.Log.Info("Emailing Service started")

	cacheService := redis.NewRedisService(
		userStore,
		productStore,
		storeStore,
		cfg,
	)
	Logger.Log.Info("Cache Service started")

	success_url := net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/success"
    cancel_url := net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/cancel"
	paymentService := paymentSrvc.NewPaymentService(cfg.STRIPE_CONF.STRIPE_SECRET_KEY, 
		cfg.STRIPE_CONF.STRIPE_WEBHOOK_SECRET,
		success_url,
		cancel_url,
	)
	Logger.Log.Info("Payment Service started")

	txManager := services.NewTxManager(db, productStore, orderStore, storeStore)

	Logger.Log.Info("Starting Server")
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
	Logger.Log.Info("Shutting Down")
	os.Exit(0)
}