package validators

import (
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hashedpassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	if err != nil {
		return err
	}
	
	return nil
}