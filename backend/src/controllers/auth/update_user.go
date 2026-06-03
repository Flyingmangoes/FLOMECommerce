package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/services"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"backend/src/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


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

func (uc *UserManager)UpdateUser() gin.HandlerFunc {
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

			newPassword = utils.PSTRING(string(hashedpass))
		}

		id := c.GetString("userId")
		email, err := uc.Users.GetEmail(c.Request.Context(), &id)
		if err != nil {
			Logger.Log.Error("Error while retrieving user email", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve email"))
			return
		}

		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &id,
				Email: email,
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
		
		hashed_password, err := uc.Users.GetPassword(c.Request.Context(), params)

		if err != nil {
    		c.Error(middlewares.ErrInternal("Failed to fetch user"))
    		return
		}

		if err := validators.ValidatePassword(*hashed_password, req.Password); err != nil {
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