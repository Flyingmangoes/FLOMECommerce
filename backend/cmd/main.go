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

	Logger.LoggerInit(cfg.Env)
	defer Logger.Log.Sync()

	db := database.NewDatabaseConnection(cfg.DBConf.DBAddr)

	userStore := repository.NewUserStore(db)
	storeStore := repository.NewStoresStore(db)
	productStore := repository.NewProductStore(db)
	orderStore := repository.NewOrderStore(db)
	tokenStore := repository.NewTokenStore(db)

	txManager := services.NewTxManager(db, productStore, orderStore, storeStore)

	serv := &server.ServerManager{
		Users: userStore,
		Stores: storeStore,
		Products: productStore,
		Orders: orderStore,
		Tokens: tokenStore,
		Tx: txManager,
		JWTSecret: []byte(cfg.ServConf.JWTSecret),
	}

	serv.Start(cfg)
}