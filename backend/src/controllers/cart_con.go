package controllers

type CartItemRequrest struct {
	ProdID 		string `json:"ProductId" binding:"required"`
	ProdName 	string `json:"ProductName" binding:"required"`
	StoreName   string `json:"storeName" binding:"required"`
	Quantity 	int    `json:"quantity" binding:"required"`
}

type AddCartRequest struct {
	Products 	[]CartItemRequrest 	`json:"items" binding:"required"`
}

type UpdateCartRequest struct {
	PIDtarget string 	`json:"productId" binding:"required"`
	NewQuantity int 	`json:"quantity" binding:"required"`
}

type DelFromCartRequest struct {
	PIDtarget string `json:"productId" binding:"required"`
}