package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	repo"backend/src/repository"
	"backend/src/services"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger"backend/src/utils/logger"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserManager struct {
    Users    	repo.UserStoreInterface
	Cart 		repo.CartStoreInterface
    Tokens   	repo.TokenStoreInterface
	JWTSecret 	[]byte
}

type UserRegisterRequest struct {
    FirstName    	string 	`json:"firstName"   binding:"required"`
    LastName     	string 	`json:"lastName"    binding:"required"`
    Username     	string 	`json:"username"    binding:"required"`

    Email       	string 	`json:"email"       binding:"required,email"`
	PhoneNumber		string 	`json:"phoneNumber" binding:"omitempty"`
    Password     	string 	`json:"password"    binding:"required,min=8"`
    UserLocale 	 	string 	`json:"userLocale"  binding:"required"`
	UserCountry  	string 	`json:"userCountry" binding:"required"`
	UserAddress  	string 	`json:"userAddress" binding:"required"`
	
    UserType     	string 	`json:"userType"    binding:"required"`
	IsAgreed	 	bool 	`json:"isAgreed" binding:"required"`
	EmailConsent  	bool	`json:"emailConsent" binding:"omitempty"`
	SmsConsent	  	bool	`json:"smsConsent" binding:"omitempty"`	
	ConsentSource 	string 	`json:"consentSrc" binding:"omitempty"`	
}

type userResponse struct {
	UserID       	string    	`json:"userId"`
    FirstName    	string    	`json:"firstName"`
    LastName    	string    	`json:"lastName"`
    Username   	  	string    	`json:"username"`
	PhoneNumber		string	  	`json:"phoneNumber"`
    Email        	string    	`json:"email"`
    UserType     	string    	`json:"userType"`
    UserLocale 	 	string    	`json:"userLocale"`
	UserCountry  	string	  	`json:"userCountry"`
	UserAddress  	string    	`json:"userAddress"`
	IsAgreed		bool 	  	`json:"isAgreed"`
    IsVerified   	bool      	`json:"isVerified"` 
    CreatedAt    	time.Time 	`json:"createdAt"`
	UpdatedAt 	 	*time.Time 	`json:"updatedAt"`
	ConsentUpdated 	*time.Time 	`json:"consentUpdated"`
}

//
// USER HANDLER IMPLEMENTATION
//

func toUserResponse(u *models.User) userResponse {
    return userResponse{
        UserID:       u.UserID,
        FirstName:    u.FirstName,
        LastName:     u.LastName,
        Username:     u.Username,
        Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
        UserType:     u.UserType,
        UserLocale:   u.Locale,
		UserCountry:  u.Country,
		UserAddress:  u.Address,
        IsVerified:   u.IsVerified,
		IsAgreed: 	  u.IsAgree,
        CreatedAt:    u.CreatedAt,
		ConsentUpdated: u.Consent_Updated,
		UpdatedAt: u.Updatedat,
    }
}

func (uc *UserManager) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserRegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request:"))
			return
		}

		hashedpass, err := services.Hashing([]byte(req.Password))
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to hash password"))
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
			HashedPassword: utils.PSTRING(string(hashedpass)),
			PhoneNumber: &req.PhoneNumber,
			UserType: &req.UserType,
		
			IsAgree: &req.IsAgreed,
			EmailConsent: &req.EmailConsent,
			SmsConsent: &req.SmsConsent,
			ConsentSource: &req.ConsentSource,
		}

		user, err := uc.Users.CreateUser(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			var PgErr *pq.Error
			if errors.As(err, &PgErr) && PgErr.Code == "23505" {
				c.Error(middlewares.ErrConflict("User already exists"))
				return
			}

			c.Error(middlewares.ErrInternal("Failed to create user"))
			return
		}

		cart, err := uc.Cart.CreateCart(c.Request.Context(), user.UserID)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("failed to create cart"))
			return
		}

		accessToken, err := jwt.GenerateAccessToken(user.UserID, user.UserType, user.IsVerified, uc.JWTSecret)
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

		Logger.Log.Info("Register process completed")
		c.JSON(http.StatusCreated, gin.H{
			"response": "user created",
			"detail": gin.H{
				"user": toUserResponse(user),
				"cart": cart,
			},

			"token": gin.H{
				"access_token": accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}