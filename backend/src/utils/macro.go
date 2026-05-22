package utils

const (
	ACTION_DELETE_USER = "delete_user"
	ACTION_DELETE_STORE = "delete_store"
	ACTION_DELETE_PRODUCT = "delete_product"
	ACTION_CANCEL_PRODUCT = "cancel_order"
	ACTION_UPDATE_STORE = "update_store"
	ACTION_UPDATE_USER = "update_user"
	ACTION_UPDATE_PRODUCT = "update_product"
)

const (
	CONFIRM_TOKEN = "confirm"
	ACCESS_TOKEN = "user"
)

const (
	DEVELOPMENT = "DEVELOPMENT"
	PRODUCTION = "PRODUCTION"
)

var AllowedAction = map[string]bool{
	ACTION_CANCEL_PRODUCT: true,
	ACTION_DELETE_PRODUCT: true,
	ACTION_DELETE_STORE: true,
	ACTION_DELETE_USER: true,
	ACTION_UPDATE_STORE: true,
	ACTION_UPDATE_USER: true,
	ACTION_UPDATE_PRODUCT: true,
}