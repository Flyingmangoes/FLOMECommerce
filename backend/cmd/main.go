package main

import (
	"backend/cmd/server"
	"backend/src/config"
	"backend/src/database"
	"backend/src/services"
	"log/slog"
	"os"

	"github.com/subosito/gotenv"
)

func main() {
	if err := gotenv.Load(); err != nil {
		slog.Warn("[WARN]", "Failed to find .env file", err)
	}

	cfg := config.NewConfig()
	if err := cfg.Validate(); err != nil {
		slog.Error("[ERROR]", "error", err)
		os.Exit(1)
	}

	db := database.NewDatabaseConnection(cfg.DBConf.DBAddr)

	userStore := services.NewUserStore(db)
	productStore := services.NewProductStore(db)
	orderStore := services.NewOrderStore(db)
	tokenStore := services.NewTokenStore(db)

	serv := server.SetupServer(userStore, productStore, orderStore, tokenStore, cfg.ServConf)

	serv.StartLoop(cfg)
}