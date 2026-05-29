package controllers

import (
	"backend/src/middlewares"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
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

		access, err := jwt.GenerateAnyToken(userId, userType, utils.SUDO_TOKEN, nil, []byte(sc.SUDOSecret))
		if err != nil {
			Logger.Log.Error("Failed to generate token", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		c.Header("X-Sudo-Token", "Sudo" + access)
		c.JSON(http.StatusOK, gin.H{"response": access})
	}
}