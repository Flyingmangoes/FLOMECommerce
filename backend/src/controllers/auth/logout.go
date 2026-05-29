package auth_controllers

import (
	"backend/src/middlewares"
	Logger "backend/src/utils/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (uc *UserManager)LogoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id string;
		var err error;

		id = c.GetString("userId")

		err = uc.Tokens.DeleteAllUserTokens(c.Request.Context(), id)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove user token"))
			return
		}

		Logger.Log.Info("Logout process completed")
		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}	
}