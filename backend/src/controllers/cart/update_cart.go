package cart

import (
	cart_types "backend/src/controllers/cart/types"
	terror "backend/src/error"
	repo_type "backend/src/repository/types"
	logger_system "backend/src/utils/logger_service"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (cm *CartManager) UpdateQuantity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req cart_types.UpdateQuantityRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.Get(c.Request.Context(), &repo_type.CartProfileParams{
			BaseParams: repo_type.BaseParams{UserId: &user_id},
		})

		if err != nil {
			logger_system.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve cart"))
			return
		}

		params := &repo_type.CartProfileParams{
			Quantity:    &req.NewQuantity,
			CartItemsID: &req.CartItemID,
			CartID:      &cart.ID,
		}

		items, err := cm.Carts.UpdateQuantity(c.Request.Context(), params)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(terror.ErrNotFound("Cart items not found"))
			}

			logger_system.Log.Error("Failed to update quantity", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to update item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": gin.H{
				"info":      "quantity updated",
				"cart":      cart,
				"cart_item": items,
			},
		})
	}
}
