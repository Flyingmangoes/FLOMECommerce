package server

import (
	"backend/src/config"
	error_service "backend/src/error"
	"backend/src/middlewares"
	repo "backend/src/repository"
	"backend/src/services"
	email_service "backend/src/services/email"
	payment_service "backend/src/services/payment"
	"backend/src/services/redis"
	logger_system "backend/src/utils/LoggerSystem"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type ServerSecret struct {
	JwtSecret 				[]byte
	SudoSecret 				[]byte
}

type ServerManager struct {
	EnvironmentStatus string

	Users 		repo.UserStoreInterface
	Products 	repo.ProductStoreInterface 
	Orders 		repo.OrderStoreInterface
	Carts 		repo.CartStoreInterface
	Tokens  	repo.TokenStoreInterface 

	Email 		*email_service.SendgridManager
	Payment 	*payment_service.PaymentService
	Tx			*services.TxManager
	Cacher		cache_service.RedisInterface

	ServerSecret
}

func (sm *ServerManager)Start(cfg *config.ConfigManager) {
	router := gin.Default()

	gin.SetMode(gin.DebugMode)
	if sm.EnvironmentStatus == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RATE_CONF.RPM), cfg.RATE_CONF.BURST)
	prison := middlewares.NewLoginPrison(cfg.APP_CONF.MAX_RETRY_LOGIN, time.Duration(cfg.APP_CONF.RETRY_LOGIN_COOLDOWN * int(time.Minute)))

	router.Use(middlewares.CORS())
	router.Use(iRate.RateLimiting())
	router.Use(error_service.JSONAppErrorReporter())

	registerRoutes(router, sm, prison)

	serverAddr := net.JoinHostPort(cfg.SERV_CONF.ServerHOST, cfg.SERV_CONF.ServerPORT)
	proxyAddr := net.JoinHostPort(cfg.SERV_CONF.ProxyHOST, cfg.SERV_CONF.ProxyPORT)
	
	router.SetTrustedProxies([]string{serverAddr, proxyAddr})

	logger_system.Log.Info("Server starting", zap.String("addr", serverAddr))
	if err := router.Run(serverAddr); err != nil {      
		logger_system.Log.Error("Server Failed to Start", zap.Error(err))
	}
}

func (sm *ServerManager)Exit() {

}