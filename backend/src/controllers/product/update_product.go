package product

import (
	product_types "backend/src/controllers/product/types"
	terror "backend/src/error"
	repo_type "backend/src/repository/types"
	logger_system "backend/src/utils/LoggerSystem"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (pm *ProductManager) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req product_types.UpdateProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("Error while binding cient request", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		params := &repo_type.ProductProfileParams{
			ProductID:  &req.ProductID,
			Name: req.NewProductName,	
			Desc: req.NewProductDesc,
			ImageUrl: req.NewImage,
			Price: req.NewPrice,
			Category: req.NewCategory,
			Availability: req.NewAvailability,
		}

		product, err := pm.Products.Update(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to update product"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"detail": gin.H{
				"info": "product updated",
				"product": product_types.CreateProductResponse(product),
			},
		})
	}
}
