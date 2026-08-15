package cart_types

type AddItemRequest struct {
	ProductID 	string 	`json:"productId" binding:"required"`
	Quantity  	int 	`json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	CartItemID 	string `json:"cartItemId" binding:"required"`
}

type UpdateQuantityRequest struct {
	CartItemID 	string	`json:"cartItemId" binding:"required"`
	NewQuantity int 	`json:"newQuantity" binding:"required"`
}