package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"

	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// PRODUCT AND STORE ONLY STRUCTURE DECLARATION
//

type StoreContext struct {
    Products 	repository.ProductStoreInterface
    Orders      repository.OrderStoreInterface
}

type ProductRequest struct {
	// Basic params required for the service
	ProductId 		string 	`json:"productId"`
	SellerId 		string 	`json:"sellerId"`
	StoreName 		string 	`json:"storeName"`
	ProductName 	string 	`json:"productName"`

	Url 			string 	`json:"url"`
	ImageUrl 		string 	`json:"imageUrl"`
	Price 			float64 `json:"price"`
	Rating 			float64 `json:"rating"`
	Desc 			string  `json:"desc"`
	Category 		string	`json:"category"`
	Availability 	int		`json:"availability"`

	// Params section required for updating 
	// a product
	NewProductName 	string 	`json:"newProductName"`
	NewProductDesc 	string	`json:"newProductDesc"`
	NewStorename 	string	`json:"newStoreName"`
	NewUrl 			string	`json:"newUrl"`
	NewImageUrl 	string 	`json:"newImageUrl"`

	NewPrice 		float64 `json:"newPrice"`
	NewCategory 	string	`json:"newCategory"`
	NewAvailability int		`json:"newAvailability"`
}

type StoreResponse struct {
	ProductId string 
	SellerId string
	StoreName string
}

//
// PRODUCT AND STORE ONLY STRUCTURE HANDLER IMPLEMENTATION
//

func (sc *StoreContext) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("(1PC) [DEBUG]", "message", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		params := &repository.ProductProfileParams{
				ProductID: req.ProductId,
				Name: req.ProductName,
				Desc: req.Desc,
				Url: req.Url,
				ImageUrl: req.ImageUrl,
				Price: req.Price,
				Category: req.Category,
				Availability: req.Availability,
		}

		product, err := sc.Products.CreateProduct(c.Request.Context(), params)
		if err != nil {
			slog.Error("(2PC) [DEBUG]", "message", err)
			c.Error(middlewares.ErrInternal("Failed to create product"))
			return
		}

		slog.Info("(0) [STATUS] Success")
		c.JSON(http.StatusCreated, gin.H{
			"message": "Product created",
			"detail": toStoreResponse(product),
		})
	}
}



func (sc *StoreContext) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}



func (sc *StoreContext) RemoveProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}



func (us *UserContext) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func toStoreResponse(p *models.Product) StoreResponse {
	return StoreResponse{

	}
}