package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	repo "backend/src/repository"
	"backend/src/services/auth"
	"backend/src/utils"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func (uc *UserManager) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user_type.UserRegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request:"))
			return
		}

		hashedpass, err := auth.Hashing([]byte(req.Password))
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to hash password"))
			return
		}

		params := &repo.UserProfileParams{
			BaseParams: repo.BaseParams{
				Email: &req.Email,
				Username: &req.Username,
				Locale: &req.UserLocale,
				Country: &req.UserCountry,
				Address: &req.UserAddress,
			},
			FirstName: &req.FirstName,
			LastName: &req.LastName,
			UserType: utils.PINT(repo.UserUnverified),
			HashedPassword: utils.PSTRING(string(hashedpass)),
			PhoneNumber: &req.PhoneNumber,
		
			IsAgree: &req.IsAgreed,
			EmailConsent: &req.EmailConsent,
			SmsConsent: &req.SmsConsent,
			ConsentSource: &req.ConsentSource,
		}

		user, err := uc.Users.CreateUser(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			var PgErr *pq.Error
			if errors.As(err, &PgErr) && PgErr.Code == "23505" {
				c.Error(terror.ErrConflict("User already exists"))
				return
			}

			c.Error(terror.ErrInternal("Failed to create user"))
			return
		}

		cart, err := uc.Cart.CreateCart(c.Request.Context(), user.UserID)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("failed to create cart"))
			return
		}

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

		logger_system.Log.Info("Register process completed")
		c.JSON(http.StatusCreated, gin.H{
			"detail": gin.H{
				"user": user_type.CreateUserResponse(user),
				"cart": cart,
			},

			"token": gin.H{
				"access_token": accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}