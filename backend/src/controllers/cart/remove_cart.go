package cart

import (
	cart_types "backend/src/controllers/cart/types"
	terror "backend/src/error"
	repo_type "backend/src/repository/types"
	logger_system "backend/src/utils/LoggerSystem"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (cm *CartManager) RemoveCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req cart_types.RemoveItemRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.Get(c.Request.Context(), &repo_type.CartProfileParams{
			BaseParams: repo_type.BaseParams{
				UserId: &user_id,
			},
		})

		if err != nil {
			logger_system.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve cart"))
			return
		}

		err = cm.Carts.RemoveItems(c.Request.Context(), &repo_type.CartProfileParams{
			CartItemsID: &req.CartItemID,
			CartID: &cart.ID,
		})

		if err != nil {
			logger_system.Log.Error("Failed to remove item", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove item"))
			return 
		}

		c.JSON(http.StatusOK, gin.H{"response": "OK"})
	}
}

func (cm *CartManager) ClearCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		requester_id := c.GetString("userId")

		cart, err := cm.Carts.Get(c.Request.Context(), &repo_type.CartProfileParams{
			BaseParams: repo_type.BaseParams{UserId: &requester_id},
		})
		err = cm.Carts.ClearItems(c.Request.Context(), &repo_type.CartProfileParams{
			BaseParams: repo_type.BaseParams{UserId:  &requester_id},
			CartID: &cart.ID,
		})
		if err != nil {
			if err == sql.ErrNoRows{
				logger_system.Log.Error("Rows not found", zap.Error(err))
				c.Error(terror.ErrNotFound("Cart not found"))
				return
			}

			logger_system.Log.Error("Error in clear cart", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to clear cart"))
			return
		}

		c.JSON(http.StatusOK, gin.H{ "response": "OK"})
	}
}