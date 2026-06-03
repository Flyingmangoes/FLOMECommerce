package middlewares

import (
	"backend/src/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func SudoMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		confirmationHeader := c.GetHeader("X-Sudo-Token")
		if confirmationHeader == "" {
			c.Error(ErrPreconditionRequired("Missing Sudo Header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(confirmationHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Sudo" {
			c.Error(ErrUnauthorized("Invalid authorization format"))
			c.Abort()
			return 
		}

		_, err := jwt.VerifySudoToken(parts[1], []byte(secret))
		if err != nil {
			c.Error(ErrUnauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Next()
	}
}