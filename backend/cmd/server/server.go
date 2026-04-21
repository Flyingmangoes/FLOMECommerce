package server

import (
	"backend/src/config"
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/utils"
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Always Store the required Interface in the server context
// to make it more readable and added the new variable
// in SetupServer

type ServerContext struct {
	Users 		repository.UserStoreInterface
	Stores 		repository.StoreStoreInterface	
	Products 	repository.ProductStoreInterface 
	Orders 		repository.OrderStoreInterface
	Tokens  	repository.TokenStoreInterface 
	JWTSecret	[]byte
}

type ServerSetupParams struct {
	Us repository.UserStoreInterface
	Ss repository.StoreStoreInterface
	Ps repository.ProductStoreInterface
	Os repository.OrderStoreInterface
	Ts repository.TokenStoreInterface
	Cfg *config.ServerConfig
}

func SetupServer(params *ServerSetupParams) *ServerContext {
	return &ServerContext{
		Users: params.Us,
		Products: params.Ps,
		Stores: params.Ss,
		Orders: params.Os,
		Tokens: params.Ts,
		JWTSecret: []byte(params.Cfg.JWTSecret),
	}
}


func (s *ServerContext)StartLoop(cfg *config.Application) {
	router := gin.Default()
	iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RateConf.RequestPerMinute), cfg.RateConf.Burst)

	router.Use(middlewares.CORSMiddleware())
	router.Use(iRate.RateLimiting())
	router.Use(middlewares.JSONAppErrorReporter())

	registerRoutes(router, s)

	addr := net.JoinHostPort(cfg.ServConf.Host, cfg.ServConf.Port)
	router.SetTrustedProxies([]string{addr})

	utils.Log.Info("Server starting", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {      
		utils.Log.Error("Server Failed to Start", zap.Error(err))
	}
}