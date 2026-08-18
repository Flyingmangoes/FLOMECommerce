package user

import (
	terror "backend/src/error"
	repo_type "backend/src/repository/types"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (uc *UserManager)DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requester_id := c.GetString("userId")
		if requester_id == "" {
			logger_system.Log.Error("Missing user id")
			c.Error(terror.ErrUnauthorized("Invalid user"))
			return
		}
		
		params := &repo_type.UserProfileParams{
			BaseParams: repo_type.BaseParams{
				UserId: &requester_id,
			},
		}

		if err := uc.Users.Delete(c.Request.Context(), params); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove user"))
			return
		}

		logger_system.Log.Info("Delete process completed")
		c.JSON(http.StatusOK, gin.H{"status": "user removed"})
	}
}