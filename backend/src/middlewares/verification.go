package middlewares

import (
	"backend/src/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerificationMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token string `json:"token"`
		}

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}

		claims, err := jwt.VerifyVerificationToken(req.Token, secret)
		if err != nil {
			c.Error(ErrUnauthorized("Invalid verification token"))
			c.Abort()
			return
		}

		c.Set("verifiedUserId", claims.UserID)
		c.Next()
	}
}

