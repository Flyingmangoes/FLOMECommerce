package cart

import (
	terror "backend/src/error"
	"backend/src/repository"
	logger_system "backend/src/utils/LoggerSystem"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (cm *CartManager) GetCarts() gin.HandlerFunc {
	return func (c *gin.Context) {
		requester_id := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repository.CartProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &requester_id,
			},
		})

		if err != nil {
			logger_system.Log.Error("Error in cart retrieval", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve cart"))
			return
		}

		items, err := cm.Carts.GetCartItems(c.Request.Context(), &repository.CartProfileParams{
			CartID: &cart.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(terror.ErrNotFound("Cart items not found"))	
				return
			}
			logger_system.Log.Error("Error in retrieving cart items", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve cart items"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": gin.H{"cart": cart,
				"cart_items":items,
			},
		})
	}
}