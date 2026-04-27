package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	"backend/src/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// PRODUCT STRUCTURE DECLARATION
//

type ProductContext struct {
	Products   repository.ProductStoreInterface
	JWTSecret []byte
}

type RegisterProductRequest struct {
	// Basic params required for the service
	StoreName 		string 	`json:"storeName"`
	ProductName 	string 	`json:"productName"`
	ImageUrl 		string 	`json:"imageUrl"`
	Price 			float64 `json:"price"`
	Rating 			float64 `json:"rating"`
	Desc 			string  `json:"desc"`
	Category 		string	`json:"category"`
	Availability 	int		`json:"availability"`
}

type UpdateProductRequest struct {
	NewProductName 	string 	`json:"newProductName"`
	NewProductDesc 	string	`json:"newProductDesc"`
	NewStorename 	string	`json:"newStoreName"`
	NewImageUrl 	string 	`json:"newImageUrl"`
	NewPrice 		float64 `json:"newPrice"`
	NewCategory 	string	`json:"newCategory"`
	NewAvailability int		`json:"newAvailability"`
}

type ProductResponse struct {
	ProductId string 
	StoreId string
	StoreName string
	StoreDesc string		
	ProductUrl string
	ProductPic string
	Price float64
	Category string
	Rating float64
	Availability int
	CreatedAt time.Time 
    UpdatedAt *time.Time  
}

//
// PRODUCT HANDLER IMPLEMENTATION
//

func (pc *ProductContext) RegisterProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			utils.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		store_id := c.GetString("storeId")

		params := &repository.ProductProfileParams{
			StoreId: &store_id,
			Name: &req.ProductName,
			ImageUrl: &req.ImageUrl,
			Price: &req.Price,
			Rating: &req.Rating,
			Desc: &req.Desc,
			Category: &req.Category,
			Availability: &req.Availability,
		}

		product, err := pc.Products.CreateProduct(c.Request.Context(), params)
		if err != nil {
			utils.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to register product"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": toProductResponse(product),
		})
	}
}



func (pc *ProductContext) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			utils.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		params := &repository.ProductProfileParams{

		}

		store, err := pc.Products.UpdateProduct(c.Request.Context(), params)
		if err != nil {
			utils.Log.Error("detail", zap.Error(err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"response": toProductResponse(store)})
	}
}



func (pc *ProductContext) RemoveProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}



func (pc *ProductContext) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func toProductResponse(p *models.Product) ProductResponse {
	return ProductResponse{
		ProductId: p.ProductID,
		StoreId: p.StoreID,
		StoreName: p.Name,
		StoreDesc: p.Desc,		
		ProductUrl: p.Url,
		ProductPic: p.ImageUrl,
		Price: p.Price,
		Category: p.Category,
		Rating: p.Rating,
		Availability: p.Availability,
		CreatedAt: p.CreatedAt,
    	UpdatedAt: p.UpdatedAt,
	}
}