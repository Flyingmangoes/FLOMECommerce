package cart

import "backend/src/repository"

type CartManager struct {
	Carts 		repository.CartStoreInterface
	Users 		repository.UserStoreInterface
	Products 	repository.ProductStoreInterface
}
