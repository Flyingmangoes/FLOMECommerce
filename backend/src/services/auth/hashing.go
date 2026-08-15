package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func Hashing(pw []byte) ([]byte, error) { 
	hashedpass, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedpass, nil
}