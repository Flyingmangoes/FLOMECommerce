package middlewares

import (
	"backend/src/repository"
	Logger "backend/src/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckForStore(stores repository.StoreStoreInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userTypes := c.GetString("userType")
		if userTypes == "" {
			c.Error(ErrBadRequest("Missing header"))
			c.Abort()
			return
		}

		if userTypes != "Store" {
			c.Error(ErrUnauthorized("Not allowed"))
			c.Abort()
			return 
		}

		ownerId := c.GetString("userId")
		store, err := stores.GetStoreByOwner(c.Request.Context(), ownerId)
		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Abort()
			return
		}
		
		c.Set("storeId", store.StoreId)
		c.Next()
	}
}