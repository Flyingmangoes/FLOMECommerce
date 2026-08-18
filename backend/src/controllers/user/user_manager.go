package user

import "backend/src/repository"

type UserManager struct {
    Users    	repository.UserStoreInterface
	Cart 		repository.CartStoreInterface
    Tokens   	repository.TokenStoreInterface
	JWTSecret 	[]byte
}