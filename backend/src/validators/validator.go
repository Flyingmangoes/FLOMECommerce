package validators

import (
	"backend/src/models"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hashedpassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	if err != nil {
		return err
	}
	
	return nil
}

func checkIsAdmin(user *models.User) bool {
	if user.UserType != "admin" {
		return false
	}
	
	return true
}