package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	"backend/src/services"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
	"backend/src/validators"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

//
// USER HANDLER STRUCTURE DECLARATION
//

type UserContext struct {
    Users    	repository.UserStoreInterface
	Cart 		repository.CartStoreInterface
    Tokens   	repository.TokenStoreInterface
	JWTSecret 	[]byte
}

type UserRegisterRequest struct {
    FirstName    string `json:"firstName"   binding:"required"`
    LastName     string `json:"lastName"    binding:"required"`
    Username     string `json:"username"    binding:"required"`

    Email        string `json:"email"       binding:"required,email"`
	PhoneNumber	 string `json:"phoneNumber" binding:"omitempty"`
    Password     string `json:"password"    binding:"required,min=8"`
    UserLocale 	 string `json:"userLocale"  binding:"required"`
	UserCountry  string `json:"userCountry" binding:"required"`
	UserAddress  string `json:"userAddress" binding:"required"`
	
    UserType     string `json:"userType"    binding:"required"`
	IsAgree		 bool 	`json:"userAgreed" binding:"required"`
	EmailConsent bool	`json:"emailConsent" binding:"omitempty"`
	SmsConsent	 bool	`json:"smsConsent" binding:"omitempty"`	
	ConsentSource string `json:"consentSrc"`	
}

type UpdateUserRequest struct {
    UserID      	string  `json:"userId"      binding:"required"`
	Password 		string 	`json:"password"`

    NewFirstname 	*string `json:"newFirstname" binding:"omitempty"`
    NewLastname  	*string `json:"newLastname"  binding:"omitempty"`
	NewUsername  	*string `json:"newUsername"  binding:"omitempty,min=3"`

	NewPhonenumber	*string	`json:"newPhonenumber" binding:"omitempty"`
	NewEmail     	*string `json:"newEmail"     binding:"omitempty,email"`
    NewPassword  	*string `json:"newPassword"  binding:"omitempty,min=8"`

    NewLocale  		*string `json:"newLocale" binding:"omitempty"`
	NewCountry		*string	`json:"newCountry" binding:"omitempty"`
	NewAddress		*string `json:"newAddress" binding:"omitempty"`

	NewEmailConsent *bool 	`json:"newEmailConsent" binding:"omitempty"`
	NewSmsConsent	*bool	`json:"newSmsConsent" binding:"omitempty"`
}

type RemoveUserRequest struct {
	UserID      string  `json:"userId"      binding:"required"`
	Email 		string  `json:"email"        binding:"required,email"`
	Password    string  `json:"password"     binding:"required"`
}

type LoginRequest struct {
	Email		 string `json:"email"    binding:"required,email"`
	Password 	 string	`json:"password" binding:"required"`
}

type userResponse struct {
	UserID       string    `json:"userId"`
    FirstName    string    `json:"firstName"`
    LastName     string    `json:"lastName"`
    Username     string    `json:"username"`
	PhoneNumber	 string	   `json:"phoneNumber"`
    Email        string    `json:"email"`
    UserType     string    `json:"userType"`
    UserLocale 	 string    `json:"userLocale"`
	UserCountry  string	   `json:"userCountry"`
	UserAddress  string    `json:"userAddress"`
	IsAgree		 bool 	   `json:"isAgree"`
    IsVerified   bool      `json:"isVerified"` 
    CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt 	 *time.Time `json:"updatedAt"`
	ConsentUpdated *time.Time `json:"consentUpdated"`
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
		IsAgree: 	  u.IsAgree,
        CreatedAt:    u.CreatedAt,
		ConsentUpdated: u.Consent_Updated,
		UpdatedAt: u.Updatedat,
    }
}

func (uc *UserContext) RegisterUser() gin.HandlerFunc {
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

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				Email: &req.Email,
				Username: &req.Username,
				Locale: &req.UserLocale,
				Country: &req.UserCountry,
				Address: &req.UserAddress,
			},
			FirstName: &req.FirstName,
			LastName: &req.LastName,
			HashedPassword: utils.PString(string(hashedpass)),
			PhoneNumber: &req.PhoneNumber,
			UserType: &req.UserType,
		
			IsAgree: &req.IsAgree,
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



func (uc *UserContext)UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		var newPassword *string = nil
		if req.NewPassword != nil {
			var pw string = *req.NewPassword

			hashedpass, err := services.Hashing([]byte(pw))
			if err != nil {
				Logger.Log.Error("Error", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to hash password"))
				return
			}

			newPassword = utils.PString(string(hashedpass))
		}

		id := c.GetString("userId")

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &id,
				Username: req.NewUsername,
				Locale: req.NewLocale,
				Country: req.NewCountry,
			},
			HashedPassword: &req.Password,

			FirstName: req.NewFirstname,
			LastName: req.NewLastname,
			NewPasswordHashed: newPassword,
			PhoneNumber: req.NewPhonenumber,
			EmailConsent: req.NewEmailConsent,
			SmsConsent: req.NewSmsConsent,
		}
		
		existingUser, err := uc.Users.GetPassword(c.Request.Context(), params)

		if err != nil {
    		c.Error(middlewares.ErrInternal("Failed to fetch user"))
    		return
		}

		if err := validators.ValidatePassword(existingUser.PasswordHash, req.Password); err != nil {
    		c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
    		return
		}

		user, err := uc.Users.UpdateUser(c.Request.Context(), params)

		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update user"))
			return
		}

		Logger.Log.Info("Update process completed")
		c.JSON(http.StatusOK, gin.H{
			"response": "user updated",
			"detail": toUserResponse(user),
		}) 
	}
}



func (uc *UserContext)DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RemoveUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		id := c.GetString("userId")

		stored, err := uc.Users.GetPassword(c.Request.Context(), &repository.UserProfileParams{
			BaseParams: repository.BaseParams{UserId: &id},
		})
		
        if err != nil {
            Logger.Log.Error("Error", zap.Error(err))
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := validators.ValidatePassword(stored.PasswordHash, req.Password); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &req.UserID,
				Email: &req.Email,
			},
		}

		if err := uc.Users.DeleteUser(c.Request.Context(), params); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove user"))
			return
		}

		Logger.Log.Info("Delete process completed")
		c.JSON(http.StatusOK, gin.H{"status": "user removed"})
	}
}



func (uc *UserContext)LoginUser(prison *middlewares.LoginPrison) gin.HandlerFunc {
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



func (uc *UserContext)Refresh() gin.HandlerFunc {
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



func (uc *UserContext)LogoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id string;
		var err error;

		id = c.GetString("userId")

		err = uc.Tokens.DeleteAllUserTokens(c.Request.Context(), id)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove user token"))
			return
		}

		Logger.Log.Info("Logout process completed")
		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}	
}