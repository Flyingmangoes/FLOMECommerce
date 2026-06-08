package middlewares

import (
	"backend/src/repository"
	"backend/src/services"
	Logger "backend/src/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StoreMiddleware(stores repository.StoreStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userTypes := c.GetString("userType")
		if userTypes == "" {
			c.Error(ErrBadRequest("Missing header"))
			c.Abort()
			return
		}

		if userTypes != string(services.AccountMerchant) {
			c.Error(ErrUnauthorized("Not allowed"))
			c.Abort()
			return 
		}

		ownerId := c.GetString("userId")
		store_id, err := stores.FetchStoreID(c.Request.Context(), ownerId)
		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Abort()
			return
		}
		
		c.Set("storeId", store_id)
		c.Next()
	}
}