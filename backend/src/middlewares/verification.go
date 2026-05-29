package middlewares

import "github.com/gin-gonic/gin"

func VerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		c.Next()
	}
}

