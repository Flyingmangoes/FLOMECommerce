package user

import (
	user_type "backend/src/controllers/user/types"
	terror "backend/src/error"
	repo_type "backend/src/repository/types"
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
			hashedpass, err := auth_service.HashPassword([]byte(*req.NewPassword))
			if err != nil {
				logger_system.Log.Error("Error", zap.Error(err))
				c.Error(terror.ErrInternal("Failed to hash password"))
				return
			}

			newPassword = utils.PSTRING(string(hashedpass))
		}

		userId := c.GetString("userId")
		params := &repo_type.UserProfileParams{
			BaseParams: repo_type.BaseParams{
				UserId: &userId,
				Username: req.NewUsername,
				PhoneNumber: req.NewPhonenumber,
			},
			FirstName: req.NewFirstname,
			LastName: req.NewLastname,
			Locale: req.NewLocale,
			Country: req.NewCountry,
			Address: req.NewAddress,
			NewPasswordHashed: newPassword,
			EmailConsent: req.NewEmailConsent,
			SmsConsent: req.NewSmsConsent,
		}

		user, err := uc.Users.Update(c.Request.Context(), params)
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