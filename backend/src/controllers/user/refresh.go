package user

import (
	terror "backend/src/error"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (uc *UserManager)Refresh() gin.HandlerFunc {
	return  func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to read request"))
			return
		}

		claims, err := jwt_service.VerifyAccessToken(req.RefreshToken, uc.JWTSecret)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrUnauthorized("Invalid or expired refresh token"))
			return
		}

		stored, err := uc.Tokens.GetRefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrUnauthorized("Refresh token not found"))
			return
		}

		if err := uc.Tokens.DeleteRefreshToken(c.Request.Context(), stored.Token); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to rotate token"))
			return
		}

		newAccess, err := jwt_service.GenerateAccessToken(claims.UserID, claims.UserType, claims.UserVerified, uc.JWTSecret)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to create token"))
			return
		}

		newRefresh, expiresAt, err := jwt_service.GenerateRefreshToken(claims.UserID, claims.UserType, uc.JWTSecret)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to create refresh token"))
			return
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), claims.UserID, newRefresh, expiresAt); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to save refresh token"))
			return
		}

		c.Header("Authorization", "Bearer" + newAccess)
		c.Header("X-Refresh-Token", newRefresh)

		logger_system.Log.Info("Refresh process completed")
		c.JSON(http.StatusOK, gin.H{
			"token": gin.H{
				"access_token": newAccess,
				"refresh_token": newRefresh,
			},
		})
	}	
}