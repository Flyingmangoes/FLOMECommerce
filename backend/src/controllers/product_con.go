package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	repo "backend/src/repository"
	"backend/src/services"
	"backend/src/services/redis"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"backend/src/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// PRODUCT STRUCTURE DECLARATION
//

type ProductManager struct {
	Stores	   repo.StoreStoreInterface
	Cache 	   redis.RedisInterface
	Products   repo.ProductStoreInterface
	JWTSecret []byte
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

func (pm *ProductManager) RegisterProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		store_id := c.GetString("storeId")

		params := &repo.ProductProfileParams{
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
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "product created",
				"product": toProductResponse(product),
			},
		})
	}
}



func (pm *ProductManager) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error while binding cient request", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		storeId := c.GetString("storeId")

		params := &repo.ProductProfileParams{
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
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "product updated",
				"product": toProductResponse(product),
			},
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

		requester_id := c.GetString("storeId")
		isAllowed, err := validators.ValidateRequester(c.Request.Context(), 
			requester_id, req.ProductID, string(services.AccountMerchant), pm.Products,
		)
		if err != nil {
			Logger.Log.Error("Error while checking credentials", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to check credentials"))
			return
		}

		if !isAllowed {
			c.Error(middlewares.ErrUnauthorized("Invalid user"))
			return
		}

		err = pm.Products.RemoveProduct(c.Request.Context(), &repo.ProductProfileParams{
			ProductID: &req.ProductID,
			StoreId: &requester_id,
		})

		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove product"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": "product removed",
		})
	}
}

func (pm *ProductManager) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SearchProductRequest

		if err := c.ShouldBindQuery(&req); err != nil {
			Logger.Log.Error("Failed to bind query", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		filter := utils.PagFilter{Limit: req.Limit}

		if req.Cursor != nil {
			cursor, err := utils.DecodeCursor(*req.Cursor)
			if err != nil {
				c.Error(middlewares.ErrBadRequest("Invalid cursor"))
				return 
			}

			filter.Cursor = cursor
		}

		filter.Normalize()

		params := &repository.ProductSearchParams{
			Query: req.Query,
			Category: req.Category,
			StoreID: req.StoreID,
			MinPrice: req.MinPrice,			
			MaxPrice: req.MaxPrice,
			SortBy: req.SortBy,
			SortOrder: req.SortOrder,
			PagFilter: filter,
		}

		url, _ := url.Parse(fmt.Sprintf("products:%s", c.Request.URL))
		cacheKey := pm.Cache.GenerateCacheKey(url)

		cached, err := pm.Cache.Get(c.Request.Context(), cacheKey)
		if err == nil && cached != nil {
			var page utils.Page[models.Product]
			if err :=  json.Unmarshal(cached, &page); err != nil {
				c.JSON(http.StatusOK, page)
				return
			}
		}

		products, err := pm.Products.SearchProduct(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to retrieve products", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve products"))
			return
		}

		page, err := utils.Build(products, filter.Limit, func(p models.Product) (time.Time, string) {
			return p.CreatedAt, p.ProductID			
		})

		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to build page"))
			return
		}

		c.JSON(http.StatusOK, page)
	}
}
