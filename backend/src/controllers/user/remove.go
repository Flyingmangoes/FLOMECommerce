package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	"backend/src/repository"
	"backend/src/services/auth"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (uc *UserManager)DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user_type.RemoveUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		requester_id := c.GetString("userId")
		password, err := uc.Users.FetchPassword(c.Request.Context(), &repository.UserProfileParams{
			BaseParams: repository.BaseParams{UserId: &requester_id},
		})
		
        if err != nil {
            logger_system.Log.Error("Error", zap.Error(err))
            c.Error(terror.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := auth.ValidatePassword(*password, req.Password); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrUnauthorized("Invalid credentials"))
			return
		}

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &requester_id,
			},
		}

		if err := uc.Users.DeleteUser(c.Request.Context(), params); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove user"))
			return
		}

		logger_system.Log.Info("Delete process completed")
		c.JSON(http.StatusOK, gin.H{"status": "user removed"})
	}
}