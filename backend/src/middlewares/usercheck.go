package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CheckForStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		userTypes := c.GetHeader("userType")
		if userTypes != "" {
			c.Error(ErrBadRequest("Missing header"))
			c.Abort()
			return
		}

		if userTypes != "Store" {
			c.Error(ErrUnauthorized("Not allowed"))
			c.Abort()
			return 
		}

		c.Next()
	}
}