package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestSanitation() gin.HandlerFunc {
	return func(c *gin.Context)	{
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	}
}

