package main

import (
	"backend/src/config"
	"backend/src/database"
	"backend/src/repository"
	"backend/src/server"
	"backend/src/services"
	paymentSrvc"backend/src/services/payment"
	Logger "backend/src/utils/logger"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load(); err != nil {
		slog.Warn(fmt.Sprintf("detail: %v", err))
	}

	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		slog.Error(fmt.Sprintf("detail: %v", err))
		os.Exit(1)
	}

	Logger.LoggerInit(cfg.ENVIRONMENT_STATUS)
	defer Logger.Log.Sync()

	db := database.NewDatabaseConnection(cfg.DB_CONF.DATABASE)

	userStore := repository.NewUserStore(db)
	storeStore := repository.NewStoresStore(db)
	productStore := repository.NewProductStore(db)
	orderStore := repository.NewOrderStore(db)
	tokenStore := repository.NewTokenStore(db)
	cartStore := repository.NewCartStore(db)

	paymentService := paymentSrvc.NewPaymentService(cfg.STRIPE_CONF.STRIPE_SECRET_KEY, 
	 	cfg.STRIPE_CONF.STRIPE_WEBHOOK_SECRET,
		net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/success",
    	net.JoinHostPort(cfg.SERV_CONF.FrontendHOST, cfg.SERV_CONF.FrontendPORT) + "/cancel", 
	)

	txManager := services.NewTxManager(db, productStore, orderStore, storeStore)

	serv := &server.ServerManager{
		Users: userStore,
		Stores: storeStore,
		Products: productStore,
		Orders: orderStore,
		Carts: cartStore,
		Tokens: tokenStore,
		Payment: paymentService,
		Tx: txManager,
		JWTSecret: []byte(cfg.SERV_CONF.JWT_SECRET),
		SUDOSecret: []byte(cfg.SERV_CONF.SUDO_SECRET),
	}

	serv.Start(cfg)
}