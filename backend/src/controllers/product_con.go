package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"	
	Logger "backend/src/utils/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// PRODUCT STRUCTURE DECLARATION
//

type ProductManager struct {
	Stores	   repository.StoreStoreInterface
	Products   repository.ProductStoreInterface
	JWTSecret []byte
}

type RegisterProductRequest struct {
	ProductName 	string 	`json:"productName"`
	ImageUrl 		string 	`json:"imageUrl"`
	Price 			float64 `json:"price"`
	Desc 			string  `json:"desc"`
	Category 		string	`json:"category"`
	Availability 	int		`json:"availability"`
}

type UpdateProductRequest struct {
	ProductID		string 		`json:"productId" binding:"required"`
	NewProductName 	*string 	`json:"newProductName" binding:"omitempty"`
	NewProductDesc 	*string		`json:"newProductDesc" binding:"omitempty"`
	NewStorename 	*string		`json:"newStoreName" binding:"omitempty"`
	NewImage 		*string 	`json:"newImage" binding:"omitempty"`
	NewPrice 		*float64 	`json:"newPrice" binding:"omitempty"`
	NewCategory 	*string		`json:"newCategory" binding:"omitempty"`
	NewAvailability *int		`json:"newAvailability" binding:"omitempty"`
	Confirmation 	bool		`json:"confirmation" binding:"required"`
}

type RemoveProductRequest struct {
	ProductID 		string	`json:"productId" binding:"required"`
	Confirmation	bool 	`json:"Confirmation" binding:"required"`
}

type productResponse struct {
	ProductId 		string 
	StoreId 		string
	StoreName 		string
	StoreDesc 		string		
	ProductUrl 		string
	ProductIMG 		string
	Price 			float64
	Category 		string
	Rating 			float64
	Availability	int
	CreatedAt 		time.Time 
    UpdatedAt 		*time.Time  
}

//
// PRODUCT HANDLER IMPLEMENTATION
//

func (pm *ProductManager) RegisterProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		store_id := c.GetString("storeId")

		params := &repository.ProductProfileParams{
			StoreId: &store_id,
			Name: &req.ProductName,
			ImageUrl: &req.ImageUrl,
			Price: &req.Price,
			Desc: &req.Desc,
			Category: &req.Category,
			Availability: &req.Availability,
		}

		product, err := pm.Products.CreateProduct(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to register product"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": "product created",
			"detail": toProductResponse(product),
		})
	}
}



func (pm *ProductManager) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		storeId := c.GetString("storeId")

		params := &repository.ProductProfileParams{
			ProductID:  &req.ProductID,
			StoreId: &storeId,
			Name: req.NewProductName,	
			Desc: req.NewProductDesc,
			ImageUrl: req.NewImage,
			Price: req.NewPrice,
			Category: req.NewCategory,
			Availability: req.NewAvailability,
		}

		product, err := pm.Products.UpdateProduct(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update product"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": "product updated",
			"detail": toProductResponse(product),
		})
	}
}



func (pm *ProductManager) RemoveProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RemoveProductRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		if req.Confirmation != true {
			c.Error(middlewares.ErrUnauthorized("Action need confirmation "))
			return
		}

		storeId := c.GetString("storeId")

		err := pm.Products.RemoveProduct(c.Request.Context(), &repository.ProductProfileParams{
			ProductID: &req.ProductID,
			StoreId: &storeId,
		})

		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove product"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "product removed"})
	}
}



func toProductResponse(p *models.Product) productResponse {
	return productResponse{
		ProductId: p.ProductID,
		StoreId: p.StoreID,
		StoreName: p.Name,
		StoreDesc: p.Desc,		
		ProductIMG: p.ProductIMG,
		Price: p.Price,
		Category: p.Category,
		Rating: p.Rating,
		Availability: p.Availability,
		CreatedAt: p.CreatedAt,
    	UpdatedAt: p.UpdatedAt,
	}
}