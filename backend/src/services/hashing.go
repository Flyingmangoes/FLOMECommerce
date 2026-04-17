package services

import (
	"backend/src/utils"

	"golang.org/x/crypto/bcrypt"
)



func Hashing(pw []byte) (*string, error) { 
	hashedpass, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return utils.Stroptr(string(hashedpass)), nil
}

func Compare() string {

}