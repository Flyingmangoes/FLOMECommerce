package server

import (
	"backend/src/config"
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/services"
	Logger "backend/src/utils/logger"
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Always Store the required Interface in the server context
// to make it more readable and added the new variable
// in SetupServer

type ServerManager struct {
	Users 		repository.UserStoreInterface
	Stores 		repository.StoreStoreInterface	
	Products 	repository.ProductStoreInterface 
	Orders 		repository.OrderStoreInterface
	Tokens  	repository.TokenStoreInterface 
	Payment 	*services.PaymentService
	Tx			*services.TxManager
	JWTSecret	[]byte
}

var RELEASE string = "RELEASE"
var DEVELOPMENT string = "DEV"

func (sm *ServerManager)Start(cfg *config.Application) {
	router := gin.Default()

	gin.SetMode(gin.DebugMode)
	if cfg.ENVIRONMENT_STATUS == RELEASE {
		gin.SetMode(gin.ReleaseMode)
	}

	iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RATE_CONF.RPM), cfg.RATE_CONF.BURST)

	router.Use(middlewares.CORSMiddleware())
	router.Use(iRate.RateLimiting())
	router.Use(middlewares.JSONAppErrorReporter())

	registerRoutes(router, sm)

	addr := net.JoinHostPort(cfg.SERV_CONF.HOST, cfg.SERV_CONF.PORT)
	accept_addr := net.JoinHostPort(cfg.SERV_CONF.ProxyHOST, cfg.SERV_CONF.ProxyPORT)
	
	router.SetTrustedProxies([]string{addr, accept_addr})

	Logger.Log.Info("Server starting", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {      
		Logger.Log.Error("Server Failed to Start", zap.Error(err))
	}
}