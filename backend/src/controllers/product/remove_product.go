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

func (pm *ProductManager) RemoveProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req product_types.RemoveProductRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		err := pm.Products.Remove(c.Request.Context(), &repo_type.ProductProfileParams{
			ProductID: &req.ProductID,
		})

		if err != nil {
			logger_system.Log.Error("detail", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to remove product"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"detail": "product removed",
		})
	}
}