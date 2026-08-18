package middlewares

import (
	terror "backend/src/error"
	"backend/src/services/auth"
	jwt_service "backend/src/utils/JWT"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddlewares(secret string) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(terror.ErrUnauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2) 
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(terror.ErrUnauthorized("Invalid authorization format"))
			c.Abort()
			return 
		}

		claims, err := jwt_service.VerifyAccessToken(parts[1], []byte(secret))
		if err != nil {
			c.Error(terror.ErrUnauthorized(err.Error()))
            c.Abort()
            return
		}

		c.Set("userVerified", claims.UserVerified)
		c.Set("userId", claims.UserID)
        c.Set("userType", claims.UserType)
        c.Next()
	}
}        

func AuthorizationMiddleware(action auth_service.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		userType := c.GetString("userType")

		ok, err := auth_service.VerifyAuthorization(userId, auth_service.AccountType(userType), action)
		if err != nil  || !ok {
			c.Error(terror.ErrUnauthorized("Invalid user"))
			c.Abort()
			return
		}

		c.Next()
	}
}