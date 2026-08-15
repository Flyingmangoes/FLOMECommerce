package order

import (
	order_types "backend/src/controllers/order/types"
	terror "backend/src/error"
	"backend/src/repository"
	order_services "backend/src/services/order"
	logger_system "backend/src/utils/LoggerSystem"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (om *OrderManager) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req order_types.OrderRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		user, err := om.Users.FetchUserByID(c.Request.Context(), &buyerId)

		if err != nil {
    		c.Error(terror.ErrInternal("Failed to fetch user"))
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
		
		order, err := om.OrderService.PlaceOrder(c.Request.Context(), &order_services.PlaceOrderParams{
			BuyerID: buyerId,
            BuyerEmail: buyerEmail,
            CombinedLocation: location,
            Status: "pending",
            ProductList: items,
		})

		if err != nil {
        	logger_system.Log.Error("Error in creating order", zap.Error(err))
            c.Error(terror.ErrInternal("Failed to place order"))
            return
		}

		c.JSON(http.StatusCreated, gin.H{
			"detail": gin.H{
				"info": "order created",
				"order": order_types.CreateOrderResponse(order),
			},
		})
	}
}
