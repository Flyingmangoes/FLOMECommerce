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

func (pm *ProductManager) RegisterProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req product_types.RegisterProductRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		params := &repo_type.ProductProfileParams{
			Name: &req.ProductName,
			ImageUrl: &req.ImageUrl,
			Price: &req.Price,
			Desc: &req.Desc,
			Category: &req.Category,
			Availability: &req.Availability,
		}

		product, err := pm.Products.Create(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to register product"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"detail": gin.H{
				"info": "product created",
				"product": product_types.CreateProductResponse(product),
			},
		})
	}
}