package middlewares

import (
	"backend/src/repository"
	"backend/src/utils"

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
		store, err := stores.GetStoreByOwner(c.Request.Context(), &repository.StoreProfileParams{OwnerId: &ownerId})
		if err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Abort()
			return
		}
		
		c.Set("storeId", store.StoreId)
		c.Next()
	}
}