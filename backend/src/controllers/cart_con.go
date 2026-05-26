package controllers

import (
	"backend/src/middlewares"
	Logger"backend/src/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AddItemRequest struct {

}

type GetCartRequest struct {

}

type RemoveItemRequest struct {
	
}

type ClearCartRequest struct {
	
}

type UpdateQuantityRequest struct {
	
}

func GetCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req GetCartRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}

func AddCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req AddItemRequest
		
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}

func UpdateQuantity() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req UpdateQuantityRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}

func RemoveCartItem() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req RemoveItemRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}

func ClearCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req ClearCartRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client requst", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}
	}
}
