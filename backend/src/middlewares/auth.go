package middlewares

import (
	"backend/src/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewares(secret string) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(ErrUnauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2) 
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(ErrUnauthorized("Invalid authorization format"))
			c.Abort()
			return 
		}

		claims, err := utils.VerifyToken(parts[1], []byte(secret))
		if err != nil {
			c.Error(ErrUnauthorized(err.Error()))
            c.Abort()
            return
		}

		c.Set("userId", claims.UserID)
        c.Set("userType", claims.UserType)
        c.Next()
	}
}