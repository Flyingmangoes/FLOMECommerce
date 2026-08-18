package auth

import "fmt"

type Action string
type AccountType string

const (
	AccountGuestUser		AccountType = "GUEST_USER"
	AccountUnverified		AccountType = "UNVERIFIED_USER"
	AccountVerified 		AccountType = "VERIFIED_USER"
	AccountAdmin 			AccountType = "ADMIN"
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

	ActionCartRead   		Action = "cart:read"
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
	AccountUnverified: {
		ActionStoreRead,
		ActionProductRead,
		ActionCartRead,
		ActionProfileDelete,
	},
	AccountVerified: {
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

func VerifyAuthorization(
		userId string, 
		userAccount AccountType, 
		userAction Action,
	) (bool, error) {
	if userAccount == "" || userAction == "" {
		return false, fmt.Errorf("Missing identifier")
	}

	actionList, exists := AuthorizationList[userAccount]
	if !exists {
		return false, fmt.Errorf("Invalid Account type: %s", string(userAccount))
	}

	isAllowed := getAllowedAction(actionList, userAction)
	if !isAllowed {
		return false, fmt.Errorf("Action not authorized for this user")
	}

	return isAllowed, nil
}

func getAllowedAction(actions []Action, userAction Action) bool {
	for _,  action := range actions {
		if userAction != action {
			return false
		}
	}

	return true
}

