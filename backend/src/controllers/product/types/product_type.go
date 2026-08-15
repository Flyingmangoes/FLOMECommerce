package product_types

import (
	"backend/src/models"
	"time"
)

type RegisterProductRequest struct {
	ProductName 	string 	`json:"productName" binding:"required"`
	ImageUrl 		string 	`json:"imageUrl" binding:"required"`
	Price 			float64 `json:"price" binding:"required"`
	Desc 			string  `json:"desc" binding:"required"`
	Category 		string	`json:"category" binding:"required"`
	Availability 	int		`json:"availability" binding:"required"`
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
}

type RemoveProductRequest struct {
	ProductID 		string	`json:"productId" binding:"required"`
}

type SearchProductRequest struct {
	Query 		*string 	`form:"q"`
	Category 	*string 	`form:"category"`
	StoreID		*string 	`form:"storeId"`
	MinPrice	*float64 	`form:"minPrice"`
	MaxPrice 	*float64 	`form:"maxPrice"`
	SortBy		*string		`form:"sortBy"`
	SortOrder 	*string		`form:"sortOrder"`
	Cursor		*string  	`form:"cursor"`
	Limit 		int			`form:"limit"`
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

func CreateProductResponse(p *models.Product) productResponse {
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