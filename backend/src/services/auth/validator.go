package auth_service

import (
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hashedPass, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		return err
	}
	
	return nil
}