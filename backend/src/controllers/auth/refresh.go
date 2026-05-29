package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
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
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to read request"))
			return
		}

		claims, err := jwt.VerifyAccessToken(req.RefreshToken, uc.JWTSecret)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid or expired refresh token"))
			return
		}

		stored, err := uc.Tokens.GetRefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Refresh token not found"))
			return
		}

		if err := uc.Tokens.DeleteRefreshToken(c.Request.Context(), stored.Token); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to rotate token"))
			return
		}

		newAccess, err := jwt.GenerateAnyToken(claims.UserID, claims.UserType, utils.ACCESS_TOKEN, nil, uc.JWTSecret)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to create token"))
			return
		}

		newRefresh, expiresAt, err := jwt.GenerateRefreshToken(claims.UserID, claims.UserType, uc.JWTSecret)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to create refresh token"))
			return
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), claims.UserID, newRefresh, expiresAt); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to save refresh token"))
			return
		}

		c.Header("Authorization", "Bearer" + newAccess)
		c.Header("X-Refresh-Token", newRefresh)

		Logger.Log.Info("Refresh process completed")
		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"token": gin.H{
				"access_token": newAccess,
				"refresh_token": newRefresh,
			},
		})
	}	
}