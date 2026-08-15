package product

import (
	repo "backend/src/repository"
	"backend/src/services/redis"
)

type ProductManager struct {
	Cache 	   cache_service.RedisInterface
	Products   repo.ProductStoreInterface
	JWTSecret []byte
}
