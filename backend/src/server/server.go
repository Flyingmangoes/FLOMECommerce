package server

import (
	"backend/src/config"
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/services"
	"backend/src/services/payment"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"net"
	"time"

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
	Carts 		repository.CartStoreInterface
	Tokens  	repository.TokenStoreInterface 
	Payment 	*payment.PaymentService
	Tx			*services.TxManager
	JWTSecret	[]byte
}

func (sm *ServerManager)Start(cfg *config.Application) {
	router := gin.Default()

	gin.SetMode(gin.DebugMode)
	if cfg.ENVIRONMENT_STATUS == utils.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RATE_CONF.RPM), cfg.RATE_CONF.BURST)
	prison := middlewares.NewLoginPrison(cfg.APP_CONF.MAX_RETRY_LOGIN, time.Duration(cfg.APP_CONF.RETRY_LOGIN_COOLDOWN * int(time.Minute)))

	router.Use(middlewares.CORSMiddleware())
	router.Use(iRate.RateLimiting())
	router.Use(middlewares.JSONAppErrorReporter())

	registerRoutes(router, sm, prison)

	addr := net.JoinHostPort(cfg.SERV_CONF.HOST, cfg.SERV_CONF.PORT)
	accept_addr := net.JoinHostPort(cfg.SERV_CONF.ProxyHOST, cfg.SERV_CONF.ProxyPORT)
	
	router.SetTrustedProxies([]string{addr, accept_addr})

	Logger.Log.Info("Server starting", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {      
		Logger.Log.Error("Server Failed to Start", zap.Error(err))
	}
}