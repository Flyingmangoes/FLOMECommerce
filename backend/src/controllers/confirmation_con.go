package controllers

import (
	terror "backend/src/error"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SudoManager struct {
	SUDOSecret string
}

func (sc *SudoManager) GenerateConfirmation() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		userType := c.GetString("userType")
		userVerified := c.GetBool("userVerified")

		if !userVerified {
			c.Error(terror.ErrUnauthorized("Invalid user"))
			return
		}

		access, err := jwt_service.GenerateSudoToken(userId, userType, []byte(sc.SUDOSecret))
		if err != nil {
			logger_system.Log.Error("Failed to generate token", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to generate token"))
			return
		}

		c.Header("X-Sudo-Token", "Sudo" + access)
		c.JSON(http.StatusOK, gin.H{
			"token": access,
		})
	}
}