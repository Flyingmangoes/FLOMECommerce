package utils

type Environment string
type Status string
type OrderStatus string

const (
	DEVELOPMENT 	Environment = "DEVELOPMENT"
	PRODUCTION 		Environment = "PRODUCTION"
)

const (
	EXIT_SUCCESS Status = "success"
	EXIT_FAILURE Status = "failed"
)

const (
	ORDER_PENDING OrderStatus = "PENDING"
	ORDER_SHIPPED OrderStatus = "SHIPPED"
	ORDER_ARRIVED OrderStatus = "ARRIVED"
)