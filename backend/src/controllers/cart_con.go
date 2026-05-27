package controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
	Logger "backend/src/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartManager struct {
	Carts repository.CartStoreInterface
	Products repository.ProductStoreInterface
}

type AddItemRequest struct {
	ProductID string `json:"productId"`
	
}

type GetCartRequest struct {

}

type RemoveItemRequest struct {
	
}

type ClearCartRequest struct {
	
}

type UpdateQuantityRequest struct {
	
}

func (cm *CartManager) GetCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req GetCartRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
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
	}
}

func (cm *CartManager) ClearCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req ClearCartRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}
