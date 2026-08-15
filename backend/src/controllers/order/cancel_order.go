package order

import (
	order_types "backend/src/controllers/order/types"
	terror "backend/src/error"
	order_services "backend/src/services/order"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (om *OrderManager) CancelOrder() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req order_types.CancelOrderRequest

		if err := c.ShouldBindBodyWithJSON(req); err != nil {
			logger_system.Log.Error("detail", zap.Error((err)))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		err := om.OrderService.CancelOrder(c.Request.Context(), &order_services.CancelOrderParams{
			BuyerId: buyerId,
			OrderId: req.OrderID,
		})

		if err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove order"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "order removed"})
	}
}