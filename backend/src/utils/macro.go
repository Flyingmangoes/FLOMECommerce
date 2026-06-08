package utils

type Action string
type TokenType int
type Environment string
type Status string
type OrderStatus string

const (
	SUDO_TOKEN 		TokenType = iota + 1
	ACCESS_TOKEN
)

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