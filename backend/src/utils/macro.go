package utils

type Action string
type TokenType int
type Environment string
type Status string
type OrderStatus string

const (
	ACTION_DELETE_USER 		Action = "delete_user"
	ACTION_DELETE_STORE 	Action = "delete_store"
	ACTION_DELETE_PRODUCT 	Action = "delete_product"
	ACTION_CANCEL_ORDER 	Action = "cancel_order"
	ACTION_UPDATE_STORE 	Action = "update_store"
	ACTION_UPDATE_USER 		Action = "update_user"
	ACTION_UPDATE_PRODUCT 	Action = "update_product"
)

const (
	SUDO_TOKEN 		TokenType = iota + 1
	ACCESS_TOKEN
)

const (
	DEVELOPMENT 	Environment = "DEVELOPMENT"
	PRODUCTION 		Environment = "PRODUCTION"
)

var AllowedAction = map[Action]bool{
	ACTION_CANCEL_ORDER: true,
	ACTION_DELETE_PRODUCT: true,
	ACTION_DELETE_STORE: true,
	ACTION_DELETE_USER: true,
	ACTION_UPDATE_STORE: true,
	ACTION_UPDATE_USER: true,
	ACTION_UPDATE_PRODUCT: true,
}

const (
	EXIT_SUCCESS Status = "success"
	EXIT_FAILURE Status = "failed"
)

const (
	ORDER_PENDING OrderStatus = "PENDING"
	ORDER_SHIPPED OrderStatus = "SHIPPED"
	ORDER_ARRIVED OrderStatus = "ARRIVED"
)