package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	"backend/src/services"
	Logger "backend/src/utils/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderManager struct {
	Orders 		repository.OrderStoreInterface
	Users 		repository.UserStoreInterface
	OrderService 	*services.OrderService
}

//
// Order Structure
//

type OrderItemRequest struct {
    ProductID 	string `json:"productId" binding:"required"`
    Quantity  	int    `json:"quantity"  binding:"required,min=1"`
}

type OrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required"`
    BuyerLocation *string       `json:"location" binding:"omitempty"`
}	

type CancelOrderRequest struct {
	OID string `json:"orderId" binding:"required"`
	Confirmation bool `json:"confirmation" binding:"required"`
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

//
// Handler
//

func toOrderResponse(o *models.Order) orderResponse {
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

/* ORDER DOCUMENTATION
* 1. CreateOrder
*	create order is a function that used by gin.IRoutes and used OrderRequest struct
*	as the parameter that then be  binded using gin.Context.ShouldBindBodyWithJSON
*
*	for buyer id aka user id we use GetString(), because the user id already set 
*	in the auth middleware which used for retrieve the user data which will be 
*	used to fill an array of OrderItemInput that used in-- 
*	OrderStoreInterface.PlaceOrder that require PlaceOrderParams to work.
*	
* 2. CheckoutOrder
* 3. CancelOrder
*/

func (om *OrderManager) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		user, err := om.Users.GetUserByID(c.Request.Context(), buyerId)

		if err != nil {
    		c.Error(middlewares.ErrInternal("Failed to fetch user"))
    		return
		}

		buyerEmail := user.Email

		var location string;

        if req.BuyerLocation == nil {
            location = fmt.Sprintf("%s, %s", user.Country, user.Address)
		}

		items := make([]repository.OrderItemInput, 0)
		for _, item := range req.Items {
			items = append(items, repository.OrderItemInput{
				ProductID: &item.ProductID,
				Quantity: &item.Quantity,
			})
		}
		
		order, err := om.OrderService.PlaceOrder(c.Request.Context(), &services.PlaceOrderParams{
			BuyerID: buyerId,
            BuyerEmail: buyerEmail,
            CombinedLocation: location,
            Status: "pending",
            ProductList: items,
		})

		if err != nil {
            Logger.Log.Error("detail", zap.Error(err))
            c.Error(middlewares.ErrInternal("Failed to place order"))
            return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": "order created",
			"detail": gin.H{
				"order": toOrderResponse(order),
			},
		})

		c.Redirect(http.StatusFound, "v1/payment/stripe")
	}
}

func (om *OrderManager) CancelOrder() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req CancelOrderRequest

		if err := c.ShouldBindBodyWithJSON(req); err != nil {
			Logger.Log.Error("detail", zap.Error((err)))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		err := om.OrderService.CancelOrder(c.Request.Context(), &services.CancelOrderParams{
			BuyerId: buyerId,
			OrderId: req.OID,
			Confirmation: req.Confirmation,
		})

		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove order"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "order removed"})
	}
}