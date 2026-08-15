package order_types

import (
	"backend/src/models"
	"time"
)

type OrderItemRequest struct {
    ProductID 	string `json:"productId" binding:"required"`
    Quantity  	int    `json:"quantity"  binding:"required,min=1"`
}

type OrderRequest struct {
	Items []OrderItemRequest 	`json:"items" binding:"required"`
    BuyerLocation *string       `json:"location" binding:"omitempty"`
}	

type CancelOrderRequest struct {
	OrderID string `json:"orderId" binding:"required"`
}

type orderResponse struct {
	ID string
	BuyerID string
	BuyerEmail string
	Total float64
	Location string
	Status string
	CreatedAt time.Time
}

func CreateOrderResponse(o *models.Order) orderResponse {
	return orderResponse{
		ID: o.ID,
		BuyerID: o.BuyerID,
		BuyerEmail: o.BuyerEmail,
		Total: o.TotalPrice,
		Location: o.Location,
		Status: o.Status,
		CreatedAt: o.CreatedAt,
	}
}


