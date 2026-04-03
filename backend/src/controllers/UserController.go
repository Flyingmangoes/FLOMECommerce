package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/services"
	"backend/src/utils"
	"backend/src/validators"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//
// USER HANDLER STRUCTURE DECLARATION
//

type UserContext struct {
    Users    	services.UserStoreInterface
    Products 	services.ProductStoreInterface
    Orders   	services.OrderStoreInterface
    Tokens   	services.TokenStoreInterface
	JWTSecret 	[]byte
}

type RegisterRequest struct {
    FirstName    string `json:"firstName"   binding:"required"`
    LastName     string `json:"lastName"    binding:"required"`
    Username     string `json:"username"    binding:"required"`

    Email        string `json:"email"       binding:"required,email"`
	PhoneNumber	 string `json:"phoneNumber" binding:"omitempty"`
    Password     string `json:"password"    binding:"required,min=8"`
    UserLocale 	 string `json:"userLocale"  binding:"required"`
	UserCountry  string `json:"userCountry" binding:"required"`
	
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

	NewEmailConsent *bool 	`json:"newEmailConsent" binding:"omitempty"`
	NewSmsConsent	*bool	`json:"newSmsConsent" binding:"omitempty"`
}

type RemoveUserRequest struct {
	UserID      string  `json:"userId"      binding:"required"`
	Email 		string  `json:"email"        binding:"required,email"`
	Password    string  `json:"password"     binding:"required"`
}

type LoginRequest struct {
	Username 	*string	`json:"username" binding:"omitempty"`
	Email		*string `json:"email"    binding:"omitempty,email"`
	Password 	 string	`json:"password" binding:"required"`
}

type UserResponse struct {
	UserID       string    `json:"userId"`
    FirstName    string    `json:"firstName"`
    LastName     string    `json:"lastName"`
    Username     string    `json:"username"`
	PhoneNumber	 string	   `json:"phoneNumber"`
    Email        string    `json:"email"`
    UserType     string    `json:"userType"`
    UserLocale 	 string    `json:"userLocale"`
	UserCountry  string	   `json:"userCountry"`
	IsAgree		 bool 	   `json:"isAgree"`
    IsVerified   bool      `json:"isVerified"` 
    CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt 	 *time.Time `json:"updatedAt"`
	ConsentUpdated *time.Time `json:"consentUpdated"`
}

//
// USER HANDLER IMPLEMENTATION
//

func (uc *UserContext)Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("(1UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request:"))
			return
		}

		hashedpass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to hash password"))
			return
		}

		params := &services.UserProfileParams{
			FirstName: &req.FirstName,
			LastName: &req.LastName,
			HashedPassword: utils.Stroptr(string(hashedpass)),
			Email: &req.Email,
			PhoneNumber: &req.PhoneNumber,

			Username: &req.Username,
			UserType: &req.UserType,
			Locale: &req.UserLocale,
			Country: &req.UserCountry,
		
			IsAgree: &req.IsAgree,
			EmailConsent: &req.EmailConsent,
			SmsConsent: &req.SmsConsent,
			ConsentSource: &req.ConsentSource,
		}

		user, err := uc.Users.CreateUser(c.Request.Context(), params)
		if err != nil {
			slog.Error("(2UC) [DEBUG]", "error", err)
			var PgErr *pq.Error
			if errors.As(err, &PgErr) && PgErr.Code == "23505" {
				c.Error(middlewares.ErrConflict("User already exists"))
				return
			}

			c.Error(middlewares.ErrInternal("Failed to create user"))
			return
		}

		accessToken, err := utils.GenerateAccessToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		refreshToken, expiresAt, err := utils.GenerateRefreshToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to generate refresh token"))
			return 
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), user.UserID, refreshToken, expiresAt); err != nil {
 	   		c.Error(middlewares.ErrInternal("Failed to save session"))
    		return
		}

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusCreated, gin.H{
			"response": toUserResponse(user),
			"token": gin.H{
				"access_token": accessToken,
				"refresh_token": refreshToken,
			},
		})
	}
}



func (uc *UserContext)Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("(1UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		var newPassword *string = nil
		if req.NewPassword != nil {
			var pw string = *req.NewPassword

			hashedpass, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			if err != nil {
				slog.Error("(2UC) [DEBUG]", "error", err)
				c.Error(middlewares.ErrInternal("Failed to hash password"))
				return
			}

			newPassword = utils.Stroptr(string(hashedpass))
		}

		id := c.GetString("userId")

		params := &services.UserProfileParams{
			UserId: &id,
			HashedPassword: &req.Password,

			FirstName: req.NewFirstname,
			LastName: req.NewLastname,
			Username: req.NewUsername,
			NewPasswordHashed: newPassword,
			PhoneNumber: req.NewPhonenumber,
			Locale: req.NewLocale,
			Country: req.NewCountry,
			EmailConsent: req.NewEmailConsent,
			SmsConsent: req.NewSmsConsent,
		}

		// Changes userid to jwt authentication later

		user, err := uc.Users.UpdateUser(c.Request.Context(), params)

		if err != nil {
			slog.Error("(3UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to update user"))
			return
		}

		if err := validators.ValidatePassword(user.PasswordHash, req.Password); err != nil {
			slog.Error("(4UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
		}

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusOK, gin.H{
			"user updated": toUserResponse(user),
		}) 
	}
}



func (uc *UserContext)Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RemoveUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("(1UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		id := c.GetString("userId")
		slog.Info("(i) [DEBUG] attempting delete", "user_id", id) // ← add this

		params := &services.UserProfileParams{
			UserId: &id,
			Email: &req.Email,
		}

		stored, err := uc.Users.GetPassword(c.Request.Context(), params)
        if err != nil {
            slog.Error("(2UC) [DEBUG]", "error", err)
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := validators.ValidatePassword(stored.PasswordHash, req.Password); err != nil {
			slog.Error("(3UC) [DEBUG]", "error", err)

			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		if err := uc.Users.DeleteUser(c.Request.Context(), params); err != nil {
			slog.Error("(4UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to remove user"))
			return
		}

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}



func (uc *UserContext)Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("(1UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		if req.Email == nil && req.Username == nil {
            c.Error(middlewares.ErrBadRequest("Email or username is required"))
            return
        }

		params:= &services.UserProfileParams{
			Email: req.Email,
		}

		user, err := uc.Users.GetUserByEmail(c.Request.Context(), params)
        if err != nil {
			slog.Error("(2UC) [DEBUG]", "error", err)
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := validators.ValidatePassword(user.PasswordHash, req.Password); err != nil {
			slog.Error("(3UC) [DEBUG]", "error", err)
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		accessToken, err := utils.GenerateAccessToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			slog.Error("(4UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		refreshToken, expiresAt, err := utils.GenerateRefreshToken(user.UserID, user.UserType, uc.JWTSecret)
		if err != nil {
			slog.Error("(5UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to generate refresh token"))
			return 
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), user.UserID, refreshToken, expiresAt); err != nil {
 	   		c.Error(middlewares.ErrInternal("Failed to save session"))
    		return
		}
		
		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusOK, gin.H{
			"response": toUserResponse(user),
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
			slog.Error("(1UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to read request"))
			return
		}

		claims, err := utils.VerifyToken(req.RefreshToken, uc.JWTSecret)
		if err != nil {
			c.Error(middlewares.ErrUnauthorized("Invalid or expired refresh token"))
			return
		}

		stored, err := uc.Tokens.GetRefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			c.Error(middlewares.ErrUnauthorized("Refresh token not found"))
			return
		}

		if err := uc.Tokens.DeleteRefreshToken(c.Request.Context(), stored.Token); err != nil {
			slog.Error("(2UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to rotate token"))
			return
		}

		newAccess, err := utils.GenerateAccessToken(claims.UserID, claims.UserType, uc.JWTSecret)
		if err != nil {
			slog.Error("(3UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to create token"))
			return
		}

		newRefresh, expiresAt, err := utils.GenerateRefreshToken(claims.UserID, claims.UserType, uc.JWTSecret)
		if err != nil {
			slog.Error("(4UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to create refresh token"))
			return
		}

		if err := uc.Tokens.SaveRefreshToken(c.Request.Context(), claims.UserID, newRefresh, expiresAt); err != nil {
			slog.Error("(5UC) [DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to save refresh token"))
			return
		}

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusOK, gin.H{
			"token": gin.H{
				"access_token": newAccess,
				"refresh_token": newRefresh,
			},
		})
	}	
}



func (uc *UserContext)Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}	
}

//
// HELPER FUNCTION SECTION
//

func toUserResponse(u *models.User) UserResponse {
    return UserResponse{
        UserID:       u.UserID,
        FirstName:    u.FirstName,
        LastName:     u.LastName,
        Username:     u.Username,
        Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
        UserType:     u.UserType,
        UserLocale:   u.Locale,
		UserCountry:  u.Country,
        IsVerified:   u.IsVerified,
		IsAgree: 	  u.IsAgree,
        CreatedAt:    u.CreatedAt,
    }
}