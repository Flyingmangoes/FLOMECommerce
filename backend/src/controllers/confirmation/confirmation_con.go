package sudo

import (
	terror "backend/src/error"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SudoManager struct {
	SudoSecret string
}

func (sm *SudoManager) CreateSudo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userId")
		userType := c.GetString("userType")
		userVerified := c.GetBool("userVerified")

		if !userVerified {
			c.Error(terror.ErrUnauthorized("Invalid user"))
			return
		}

		sudoToken, err := jwt_service.GenerateSudoToken(userId, userType, []byte(sm.SudoSecret))
		if err != nil {
			logger_system.Log.Error("Failed to generate token", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to generate token"))
			return
		}

		c.Header("X-Sudo-Token", "Sudo" + sudoToken)
		c.JSON(http.StatusOK, gin.H{
			"token": sudoToken,
		})
	}
}