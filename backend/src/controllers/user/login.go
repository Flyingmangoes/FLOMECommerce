package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/services/auth"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Email		 string `json:"email"    binding:"omitempty,email"`
	Username 	 string `json:"username" binding:"omitempty"`
	Password 	 string	`json:"password" binding:"required,min=8"`
}

func (uc *UserManager)LoginUser(prison *middlewares.LoginPrison) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		var req LoginRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		user, err := uc.Users.LoginUser(c.Request.Context(), &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				Username: &req.Username,
				Email: &req.Email,
			},
		})
		if err != nil {
			logger_system.Log.Error("Error while retrieving user", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to fetch user data"))
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

		if err := auth.ValidatePassword(user.PasswordHash, req.Password); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			prison.RecordFailure(key)
            c.Error(terror.ErrUnauthorized("Invalid credentials"))
            return
        }

		prison.Release(key)

		accessToken, err := jwt_service.GenerateAccessToken(user.UserID, user.UserType, user.IsVerified, uc.JWTSecret)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to generate token"))
			return
		}

		refreshToken, expiresAt, err := jwt_service.GenerateRefreshToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to generate refresh token"))
			return 
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), user.UserID, refreshToken, expiresAt); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to save session"))
    		return
		}

		c.Header("Authorization", "Bearer" + accessToken)
		c.Header("X-Refresh-Token", refreshToken)

		logger_system.Log.Info("Login process completed")
		c.JSON(http.StatusOK, gin.H{
			"response": "login success",
			"detail": user_type.CreateUserResponse(user),

			// for postman remove after make frontend
			"token": gin.H{
				"access_token": accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}