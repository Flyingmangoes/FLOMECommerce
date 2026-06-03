package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
	Logger"backend/src/utils/logger"
	"backend/src/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RemoveUserRequest struct {
	UserID      string  `json:"userId"      binding:"required"`
	Email 		string  `json:"email"        binding:"required,email"`
	Password    string  `json:"password"     binding:"required"`
}

func (uc *UserManager)DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RemoveUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		id := c.GetString("userId")

		password, err := uc.Users.GetPassword(c.Request.Context(), &repository.UserProfileParams{
			BaseParams: repository.BaseParams{UserId: &id},
		})
		
        if err != nil {
            Logger.Log.Error("Error", zap.Error(err))
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := validators.ValidatePassword(*password, req.Password); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &req.UserID,
				Email: &req.Email,
			},
		}

		if err := uc.Users.DeleteUser(c.Request.Context(), params); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove user"))
			return
		}

		Logger.Log.Info("Delete process completed")
		c.JSON(http.StatusOK, gin.H{"status": "user removed"})
	}
}