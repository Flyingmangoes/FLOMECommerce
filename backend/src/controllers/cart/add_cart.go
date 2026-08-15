package cart

import (
	cart_types "backend/src/controllers/cart/types"
	terror "backend/src/error"
	"backend/src/repository"
	logger_system "backend/src/utils/LoggerSystem"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (cm *CartManager) AddCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req cart_types.AddItemRequest
		
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repository.CartProfileParams{
			BaseParams: repository.BaseParams{UserId: &user_id},
		})

		if err == sql.ErrNoRows {
			cart, err = cm.Carts.CreateCart(c.Request.Context(), user_id)
			if err != nil {
				logger_system.Log.Error("Failed to create cart", zap.Error(err))
				c.Error(terror.ErrInternal("Failed to create cart"))
				return
			}
		} else if err != nil {
			logger_system.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve cart"))
			return
		}

		store_id, err := cm.Products.FetchStoreID(c.Request.Context(), req.ProductID)
		if err != nil {
			logger_system.Log.Error("Failed to retrieve product", zap.Error(err))
			c.Error(terror.ErrInternal("Faile to retrieve product"))
			return
		}

		params := &repository.CartProfileParams{
			BaseParams: repository.BaseParams{UserId: &user_id},
			CartID: 	&cart.ID,
			ProductID: 	&req.ProductID,
			StoreID: &store_id,
			Quantity: 	&req.Quantity,
		}

		cart_items, err := cm.Carts.AddCartItem(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("Failed to add item", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to add item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": gin.H{
				"info":"item added",
				"cart": cart,
				"cart_item": cart_items,
			},
		})
	}
}
