package middlewares

import (
	"backend/src/repository"
	"github.com/gin-gonic/gin"
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
		
		c.Set("storeId", store.StoreId)
		c.Next()
	}
}