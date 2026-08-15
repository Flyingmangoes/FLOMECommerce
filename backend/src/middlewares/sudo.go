package middlewares

import (
	terror "backend/src/error"
	jwt_service "backend/src/utils/JWT"
	"strings"

	"github.com/gin-gonic/gin"
)

func SudoMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		confirmationHeader := c.GetHeader("X-Sudo-Token")
		if confirmationHeader == "" {
			c.Error(terror.ErrPreconditionRequired("Missing Sudo Header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(confirmationHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Sudo" {
			c.Error(terror.ErrUnauthorized("Invalid authorization format"))
			c.Abort()
			return 
		}

		_, err := jwt_service.VerifySudoToken(parts[1], []byte(secret))
		if err != nil {
			c.Error(terror.ErrUnauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Next()
	}
}