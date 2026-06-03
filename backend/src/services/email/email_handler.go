package email_services

import (
	"backend/src/middlewares"
	"backend/src/utils"
	"backend/src/utils/jwt"
	Logger "backend/src/utils/logger"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func(sg *SGMailManager) SendVerificationMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetString("userId")

		user, err := sg.Users.GetUserByID(c.Request.Context(), &user_id)
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

		v_token, expires, err := jwt.GenerateVerificationToken(user_id, []byte(sg.VERIFICATION_SECRET))
		if err != nil {
			Logger.Log.Error("Error in generating token", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to generate token"))
			return
		}

		params := &MailServiceParams{
			From: sg.TEST_EMAIL,
			To: []string{user.Email},
			Subject: "Email Verification for Flommerce Account",
			MailType: MailConfirmation,
			MailData: &MailData{
				Username: user.Username,
				UserEmail: user.Email,
				Token: v_token,
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
			"detail": user.Email,
		})
	}
}

func(sg *SGMailManager) VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		verified_id := c.GetString("verifiedUserId")

		user, err := sg.Users.VerifyUser(c.Request.Context(), verified_id)
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