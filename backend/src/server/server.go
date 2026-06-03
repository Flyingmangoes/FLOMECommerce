package server

import (
	"backend/src/config"
	"backend/src/middlewares"
	repo"backend/src/repository"
	"backend/src/services"
	emailSrvc"backend/src/services/email"
	paymentSrvc"backend/src/services/payment"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type ServerManager struct {
	Users 		repo.UserStoreInterface
	Stores 		repo.StoreStoreInterface	
	Products 	repo.ProductStoreInterface 
	Orders 		repo.OrderStoreInterface
	Carts 		repo.CartStoreInterface
	Tokens  	repo.TokenStoreInterface 

	Email 		*emailSrvc.SGMailManager
	Payment 	*paymentSrvc.PaymentService
	Tx			*services.TxManager

	JWTSecret			[]byte
	SUDOSecret 			[]byte
	VERIFICATION_SECRET []byte
}

func (sm *ServerManager)Start(cfg *config.ConfigManager) {
	router := gin.Default()

	gin.SetMode(gin.DebugMode)
	if utils.Environment(cfg.ENVIRONMENT_STATUS) == utils.PRODUCTION {
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