package main

import (
	"backend/cmd/server"
	"backend/src/config"
	"backend/src/database"
	"backend/src/repository"
	"backend/src/utils"
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

	utils.LoggerInit(cfg.Env)
	defer utils.Log.Sync()

	db := database.NewDatabaseConnection(cfg.DBConf.DBAddr)

	userStore := repository.NewUserStore(db)
	storeStore := repository.NewStoresStore(db)
	productStore := repository.NewProductStore(db)
	orderStore := repository.NewOrderStore(db)
	tokenStore := repository.NewTokenStore(db)

	params := &server.ServerSetupParams{
		Us: userStore,
		Ss: storeStore,
		Ps: productStore,
		Os: orderStore,
		Ts: tokenStore,
		Cfg: cfg.ServConf,
	}

	serv := server.SetupServer(params)

	serv.StartLoop(cfg)
}