package auth_service

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) { 
	hashedpass, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedpass, nil
}