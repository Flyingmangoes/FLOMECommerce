package email_services

import (
	error_service "backend/src/error"
	jwt_service "backend/src/utils/JWT"
	logger_system "backend/src/utils/LoggerSystem"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func(sg *SendgridManager) UserVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetString("userId")

		user, err := sg.Users.Fetch(c.Request.Context(), &user_id)
		if err != nil {
			if err == sql.ErrNoRows {
				logger_system.Log.Error("Data not found")
				c.Status(http.StatusNotFound)
				return
			}
			logger_system.Log.Error("Error while getting user data", zap.Error(err))
			c.Error(error_service.ErrInternal("Failed to retrieve user data"))
			return
		}

		token, expires, err := jwt_service.GenerateUserVerificationToken(user_id, []byte(sg.UserVerificationSecret))
		if err != nil {
			logger_system.Log.Error("Error in generating token", zap.Error(err))
			c.Error(error_service.ErrInternal("Failed to generate token"))
			return
		}

		token_url := fmt.Sprintf("http://%s/v2/user/verify?token=%s", sg.ServerUrl, token)

		params := &MailServiceParams{
			From: sg.TestEmail,
			To: []string{user.Email},
			Subject: "Verified your Account",
			MailType: UserVerification,
			MailData: &MailData{
				DomainName: sg.DomainEmail,
				SupportEmail: sg.CustomerSupportEmail,
				FirstName: user.FirstName,
				Username: user.Username,
				TokenUrl: token_url,
				Expiration: &expires,
			},
		}

		mailReq := sg.NewMail(params)
		if err := sg.SendMail(mailReq); err != nil {
			logger_system.Log.Error("Failed to send email", zap.Error(err))
			c.Error(error_service.ErrInternal("Failed to send email"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": gin.H{
				"info": "Sended",
				"email": user.Email,
			},
		})
	}
}

func(sg *SendgridManager) VerifyUserVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token string `form:"token"`
		}

		if err := c.ShouldBindQuery(&req); err != nil {
			logger_system.Log.Error("Failed to bind query", zap.Error(err))
			c.Error(error_service.ErrBadRequest("Failed to read client request"))
			return
		}

		claims, err := jwt_service.VerifyUserVerificationToken(req.Token, []byte(sg.UserVerificationSecret))
		if err != nil {
			c.Error(error_service.ErrUnauthorized("Invalid or Expired token"))
			return
		}

		user, err := sg.Users.VerifyUser(c.Request.Context(), claims.UserID)
		if err != nil {
			logger_system.Log.Error("Failed to verified user", zap.Error(err))
			c.Error(error_service.ErrInternal("Failed to verified user"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": gin.H{
				"info":"Email verified",
				"id": user.UserID,
				"email": user.Email,
			},
		})
	}
}

func (sg *SendgridManager) SendPassReset() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SendgridManager) VerifiyPassReset() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SendgridManager) SendOrderConfirmation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SendgridManager) VerifyOrderConfirmation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}