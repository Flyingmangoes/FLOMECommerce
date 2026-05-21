package main

import (
	"backend/cmd/server"
	"backend/src/config"
	"backend/src/database"
	"backend/src/repository"
	"backend/src/services"
	Logger "backend/src/utils/logger"
	"log/slog"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load(); err != nil {
		slog.Warn("DEBUG", "detail", err)
	}

	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		slog.Error("DEBUG", "detail", err)
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
	payment := services.SetupPayment(cfg)

	txManager := services.NewTxManager(db, productStore, orderStore, storeStore)

	serv := &server.ServerManager{
		Users: userStore,
		Stores: storeStore,
		Products: productStore,
		Orders: orderStore,
		Tokens: tokenStore,
		Payment: payment,
		Tx: txManager,
		JWTSecret: []byte(cfg.SERV_CONF.JWT_SECRET),
	}

	serv.Start(cfg)
}