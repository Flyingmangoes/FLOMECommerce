package order

import (
	"backend/src/repository"
	"backend/src/services/order"
)

type OrderManager struct {
	Orders 			repository.OrderStoreInterface
	Users 			repository.UserStoreInterface
	OrderService 	*order_services.OrderService
}