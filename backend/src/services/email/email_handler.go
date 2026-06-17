package email_services

import (
	"backend/src/middlewares"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func(sg *SGMailManager) SendUserVerificationMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetString("userId")

		user, err := sg.Users.FetchUserByID(c.Request.Context(), &user_id)
		if err != nil {
			if err == sql.ErrNoRows {
				Logger.Log.Error("Data not found")
				c.Status(http.StatusNotFound)
				return
			}
			Logger.Log.Error("Error while getting user data", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve user data"))
			return
		}

		token, expires, err := jwt.GenerateUserVerificationToken(user_id, []byte(sg.USER_VERIFICATION_SECRET))
		if err != nil {
			Logger.Log.Error("Error in generating token", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		token_url := fmt.Sprintf("http://%s/v2/user/verify?token=%s", sg.SERVER_URL, token)

		params := &MailServiceParams{
			From: sg.TEST_EMAIL,
			To: []string{user.Email},
			Subject: "Verified your Account",
			MailType: UserVerification,
			MailData: &MailData{
				DomainName: "Flommerce",
				SupportEmail: "support@flommerce.com",
				FirstName: user.FirstName,
				Username: user.Username,
				TokenUrl: token_url,
				Expiration: &expires,
			},
		}

		mailReq := sg.NewMail(params)
		if err := sg.SendMail(mailReq); err != nil {
			Logger.Log.Error("Failed to send email", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to send email"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "Sended",
				"email": user.Email,
			},
		})
	}
}

func(sg *SGMailManager) VerifyUserVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token string `form:"token"`
		}

		if err := c.ShouldBindQuery(&req); err != nil {
			Logger.Log.Error("Failed to bind query", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		claims, err := jwt.VerifyUserVerificationToken(req.Token, []byte(sg.USER_VERIFICATION_SECRET))
		if err != nil {
			c.Error(middlewares.ErrUnauthorized("Invalid or Expired token"))
			return
		}

		user, err := sg.Users.VerifyUser(c.Request.Context(), claims.UserID)
		if err != nil {
			Logger.Log.Error("Failed to verified user", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to verified user"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info":"Email verified successfully",
				"id": user.UserID,
				"email": user.Email,
			},
		})
	}
}

func (sg *SGMailManager) SendPassResetMail() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SGMailManager) VerifiyPassReset() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SGMailManager) SendOrderConfirmation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (sg *SGMailManager) VerifyOrderConfirmation() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}