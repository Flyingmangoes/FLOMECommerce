package services

import(
	Logger"backend/src/utils/logger"
	"go.uber.org/zap"
	"fmt"
)

type Action string
type AccountType string

const (
	AccountGuestUser	AccountType = "GUEST_USER"
	AccountNormalUser 	AccountType = "NORMAL_USER"
	AccountMerchant 	AccountType = "MERCHANT"
	AccountAdmin 		AccountType = "ADMIN"
)

const (
	ActionProfileUpdate  	Action = "profile:update"
	ActionProfileDelete  	Action = "profile:delete"
	ActionProfileSelfRead 	Action = "profile:read"

	ActionProductRead 	 	Action = "product:read"
	ActionProductCreate 	Action = "product:create"
	ActionProductUpdate 	Action = "product:update"
	ActionProductDelete 	Action = "product:delete"

	ActionStoreRead 	 	Action = "store:read"
	ActionStoreUpdate 	 	Action = "store:update"
	ActionStoreDelete	 	Action = "store:delete"
	
	ActionOrderCreate 	 	Action = "order:create"
	ActionOrderCancel 	 	Action = "order:cancel"
	ActionOrderDelete 	 	Action = "order:delete"
	ActionOrderSelfRead  	Action = "order:read"

	ActionCartSelfRead   	Action = "cart:read"
	ActionCartAdd 		 	Action = "cart:add"
	ActionCartUpdate 	 	Action = "cart:update"
	ActionCartClear 	 	Action = "cart:clear"
	ActionCartRemove 	 	Action = "cart:remove"
)

var AuthorizationList = map[AccountType][]Action{
	AccountGuestUser: {
		ActionProductRead,
		ActionStoreRead,
	},
	AccountNormalUser: {
		ActionProfileUpdate,
		ActionProfileDelete,
		ActionProductRead,
		ActionStoreRead,
		ActionOrderCreate,
		ActionOrderCancel,
		ActionCartAdd,
		ActionCartRemove,
		ActionCartClear,
		ActionCartUpdate,
	},
	AccountMerchant: {
		ActionProductCreate,
		ActionProductDelete,
		ActionProductUpdate,
		ActionProductRead,
		ActionStoreRead,
		ActionStoreUpdate,
		ActionStoreDelete,
	},
	AccountAdmin: {
		ActionProfileUpdate, 
		ActionProductRead, 	
		ActionProductCreate, 
		ActionProductUpdate, 
		ActionProductDelete,
		ActionStoreRead, 
		ActionStoreUpdate, 	
		ActionOrderCreate, 	
		ActionOrderCancel,
		ActionOrderDelete,
		ActionCartAdd, 		
		ActionCartUpdate, 	
		ActionCartClear, 	
		ActionCartRemove, 	
	},
}

const (
	allowed bool = true
	notAllowed bool = false
)

func VerifyAuthorization(user_account AccountType, user_action Action) (bool, error) {
	if user_account == "" || user_action == "" {
		return false, fmt.Errorf("Any of the parameter cannot be empty")
	}

	action_list, exists := AuthorizationList[user_account]
	if !exists {
		Logger.Log.Error("Unknown account type", zap.String("type", string(user_account)))
		return false, fmt.Errorf("Unknown Account type: %s", string(user_account))
	}

	if !getAllowedAction(action_list, user_action){
		return notAllowed, fmt.Errorf("Action not authorized for this user")
	}

	return allowed, nil
}

func getAllowedAction(actions []Action, user_action Action) bool {
	for _,  allowed_action := range actions {
		if user_action != allowed_action {
			return notAllowed
		}
	}

	return allowed
}

