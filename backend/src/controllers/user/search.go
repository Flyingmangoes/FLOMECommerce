package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func(um *UserManager) SearchUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user_type.SearchUserRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Failed to parse client request", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to parse client request"))
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}