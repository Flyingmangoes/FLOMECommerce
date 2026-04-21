package controllers

import (
	"backend/src/models"
	"backend/src/repository"

	"github.com/gin-gonic/gin"
)

//
// PRODUCT AND STORE ONLY STRUCTURE DECLARATION
//

type ProductContext struct {
	Products   repository.ProductStoreInterface
	JWTSecret []byte
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