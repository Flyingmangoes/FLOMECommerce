package controllers

import (
	"backend/src/middlewares"
	repo "backend/src/repository"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CartManager struct {
	Carts 		repo.CartStoreInterface
	Products 	repo.ProductStoreInterface
}

type AddItemRequest struct {
	ProductID 	string 	`json:"productId" binding:"required"`
	Quantity  	int 	`json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	CartItemId 	string `json:"cartItemId" binding:"required"`
}

type ClearCartRequest struct {
	CartId string `json:"cartId" binding:"required"`
}

type UpdateQuantityRequest struct {
	CartItemID 	string	`json:"cartItemId" binding:"required"`
	ProductID 	string 	`json:"productId" binding:"required"`
	NewQuantity int 	`json:"newQuantity" binding:"required"`
}

func (cm *CartManager) GetCarts() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req struct {
			CartId string `json:"cartId" binding:"required"`
		}

		userId := c.GetString("userId")
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repo.CartProfileParams{
			CartID: &userId,
		})

		if err != nil {
			Logger.Log.Error("Error in cart retrieval", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		req.CartId = cart.ID
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error in binding request", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		items, err := cm.Carts.GetCartItems(c.Request.Context(), &repo.CartProfileParams{
			CartID: &req.CartId,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(middlewares.ErrNotFound("Cart items not found"))	
			}
			Logger.Log.Error("Error in retrieving cart items", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart items"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response":utils.EXIT_SUCCESS,
			"detail": gin.H{"cart": cart,
				"cart_items":items,
			},
		})
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
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repo.CartProfileParams{
			BaseParams: repo.BaseParams{UserId: &user_id},
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

		store_id, err := cm.Products.FetchStoreID(c.Request.Context(), req.ProductID)
		if err != nil {
			Logger.Log.Error("Failed to retrieve product", zap.Error(err))
			c.Error(middlewares.ErrInternal("Faile to retrieve product"))
			return
		}

		params := &repo.CartProfileParams{
			BaseParams: repo.BaseParams{UserId: &user_id},
			CartID: 	&cart.ID,
			ProductID: 	&req.ProductID,
			StoreID: &store_id,
			Quantity: 	&req.Quantity,
		}

		cart_items, err := cm.Carts.AddCartItem(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to add item", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to add item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info":"item added",
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
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repo.CartProfileParams{
			BaseParams: repo.BaseParams{UserId: &user_id},
		})

		if err != nil {
			Logger.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		params := &repo.CartProfileParams{
			Quantity: &req.NewQuantity,
			CartItemsID: &req.CartItemID,
			CartID: &cart.ID,
		}

		items, err := cm.Carts.UpdateItemQuantity(c.Request.Context(), params)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Error(middlewares.ErrNotFound("Cart items not found"))	
			}

			Logger.Log.Error("Failed to update quantity", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update item"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info":"quantity updated",
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
		cart, err := cm.Carts.GetCart(c.Request.Context(), &repo.CartProfileParams{
			BaseParams: repo.BaseParams{
				UserId: &user_id,
			},
		})

		if err != nil {
			Logger.Log.Error("Failed to retrieve cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve cart"))
			return
		}

		err = cm.Carts.RemoveItem(c.Request.Context(), &repo.CartProfileParams{
			CartItemsID: &req.CartItemId,
			CartID: &cart.ID,
		})

		if err != nil {
			Logger.Log.Error("Failed to remove item", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove item"))
			return 
		}

		c.JSON(http.StatusOK, gin.H{"response": utils.EXIT_SUCCESS})
	}
}

func (cm *CartManager) ClearCart() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req ClearCartRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error in reading request", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		err := cm.Carts.ClearCart(c.Request.Context(), &repo.CartProfileParams{CartID: &req.CartId})
		if err != nil {
			if err == sql.ErrNoRows{
				Logger.Log.Error("Rows not found", zap.Error(err))
				c.Error(middlewares.ErrNotFound("Cart not found"))
				return
			}

			Logger.Log.Error("Error in clear cart", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to clear cart"))
			return
		}

		c.JSON(http.StatusOK, gin.H{ "response": utils.EXIT_SUCCESS})
	}
}
