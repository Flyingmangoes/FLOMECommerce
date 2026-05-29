package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
	"backend/src/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Email		 string `json:"email"    binding:"required,email"`
	Password 	 string	`json:"password" binding:"required"`
}

func (uc *UserManager)LoginUser(prison *middlewares.LoginPrison) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		var req LoginRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}


		user, err := uc.Users.LoginByUserEmail(c.Request.Context(), &req.Email)
        if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		locked, remaining := prison.IsLocked(key)
		if locked {
			c.JSON(http.StatusTooManyRequests, gin.H{
                "error":      "too many failed attempts",
                "retry_after": remaining.String(),
            })
            return
		}

		if err := validators.ValidatePassword(user.PasswordHash, req.Password); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			prison.RecordFailure(key)
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		prison.Release(key)

		accessToken, err := jwt.GenerateAnyToken(user.UserID, user.UserType, utils.ACCESS_TOKEN, nil, uc.JWTSecret)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		refreshToken, expiresAt, err := jwt.GenerateRefreshToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to generate refresh token"))
			return 
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), user.UserID, refreshToken, expiresAt); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to save session"))
    		return
		}

		c.Header("Authorization", "Bearer" + accessToken)
		c.Header("X-Refresh-Token", refreshToken)

		Logger.Log.Info("Login process completed")
		c.JSON(http.StatusOK, gin.H{
			"response": "login success",
			"detail": toUserResponse(user),

			// for postman remove after make frontend
			"token": gin.H{
				"access_token": accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}