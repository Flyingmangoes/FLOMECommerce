package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	"backend/src/repository"
	"backend/src/services/auth"
	"backend/src/utils"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (uc *UserManager)UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user_type.UpdateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		var newPassword *string = nil
		if req.NewPassword != nil {
			hashedpass, err := auth.Hashing([]byte(*req.NewPassword))
			if err != nil {
				logger_system.Log.Error("Error", zap.Error(err))
				c.Error(terror.ErrInternal("Failed to hash password"))
				return
			}

			newPassword = utils.PSTRING(string(hashedpass))
		}

		userId := c.GetString("userId")
		params := &repository.UserProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &userId,
				Username: req.NewUsername,
				Locale: req.NewLocale,
				Country: req.NewCountry,
			},
			
			FirstName: req.NewFirstname,
			LastName: req.NewLastname,
			NewPasswordHashed: newPassword,
			PhoneNumber: req.NewPhonenumber,
			EmailConsent: req.NewEmailConsent,
			SmsConsent: req.NewSmsConsent,
		}
		
		hashed_password, err := uc.Users.FetchPassword(c.Request.Context(), params)
		if err != nil {
    		c.Error(terror.ErrInternal("Failed to fetch user"))
    		return
		}

		if err := auth.ValidatePassword(*hashed_password, req.Password); err != nil {
    		c.Error(terror.ErrUnauthorized("Invalid credentials"))
    		return
		}

		user, err := uc.Users.UpdateUser(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("Error", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to update user"))
			return
		}

		logger_system.Log.Info("Update process completed")
		c.JSON(http.StatusOK, gin.H{
			"detail": user_type.CreateUserResponse(user),
		}) 
	}
}