package user

import (
	terror "backend/src/error"
	logger_system "backend/src/utils/LoggerSystem"
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
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove user token"))
			return
		}

		logger_system.Log.Info("Logout process completed")
		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}	
}