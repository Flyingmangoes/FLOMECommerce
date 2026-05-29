package controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
	Logger "backend/src/utils/logger"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartManager struct {
	Carts repository.CartStoreInterface
	Products repository.ProductStoreInterface
}

type AddItemRequest struct {
	ProductID 	string 	`json:"productId" binding:"required"`
	Quantity  	int 	`json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	CartItemId 	string `json:"cartItemId" binding:"required"`
}

type UpdateQuantityRequest struct {
	CartItemID 	string	`json:"cartItemId" binding:"required"`
	ProductID 	string 	`json:"productId" binding:"required"`
	NewQuantity int 	`json:"newQuantity" binding:"required"`
}

func (cm *CartManager) GetCarts() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func (cm *CartManager) AddCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req AddItemRequest
		
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repository.CartProfileParams{
			BaseParams: repository.BaseParams{UserId: &user_id},
		})

		if err == sql.ErrNoRows {
			cart, err = cm.Carts.CreateCart(c.Request.Context(), user_id)
			if err != nil {
				Logger.Log.Error("Failed to create cart", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to create cart"))
				return
			}
		} else if err != nil {
			Logger.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		product, err := cm.Products.GetProductByID(c.Request.Context(), req.ProductID)
		if err != nil {
			Logger.Log.Error("Failed to retrieve product", zap.Error(err))
			c.Error(middlewares.ErrInternal("Faile to retrieve product"))
			return
		}

		params := &repository.CartProfileParams{
			BaseParams: repository.BaseParams{UserId: &user_id},
			CartID: 	&cart.ID,
			ProductID: 	&req.ProductID,
			StoreID: &product.StoreID,
			Quantity: 	&req.Quantity,
		}

		cart_items, err := cm.Carts.AddCartItem(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to add item", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to add item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"detail": gin.H{
				"cart": cart,
				"cart_item": cart_items,
			},
		})
	}
}

func (cm *CartManager) UpdateQuantity() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req UpdateQuantityRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repository.CartProfileParams{
			BaseParams: repository.BaseParams{UserId: &user_id},
		})

		if err != nil {
			Logger.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		params := &repository.CartProfileParams{
			Quantity: &req.NewQuantity,
			CartItemsID: &req.CartItemID,
			CartID: &cart.ID,
		}

		items, err := cm.Carts.UpdateItemQuantity(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to update quantity", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": "success",
			"detail": gin.H{
				"cart": cart,
				"cart_item": items,
			},
		})
	}
}

func (cm *CartManager) RemoveCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req RemoveItemRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		user_id := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repository.CartProfileParams{
			BaseParams: repository.BaseParams{
				UserId: &user_id,
			},
		})

		if err != nil {
			Logger.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		err = cm.Carts.RemoveItem(c.Request.Context(), &repository.CartProfileParams{
			CartItemsID: &req.CartItemId,
			CartID: &cart.ID,
		})

		if err != nil {
			Logger.Log.Error("Failed to remove item", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove item"))
			return 
		}

		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}
}

func (cm *CartManager) ClearCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"response": "success"})
	}
}
