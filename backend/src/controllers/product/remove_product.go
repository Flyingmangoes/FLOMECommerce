package product

import (
	product_types "backend/src/controllers/product/types"
	terror "backend/src/error"
	"backend/src/repository"
	"backend/src/services/auth"
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

		requester_id := c.GetString("storeId")
		isAllowed, err := auth.ValidateRequester(c.Request.Context(), 
			requester_id, req.ProductID, string(auth.AccountAdmin), pm.Products,
		)
		if err != nil {
			logger_system.Log.Error("Error while checking credentials", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to check credentials"))
			return
		}

		if !isAllowed {
			c.Error(terror.ErrUnauthorized("Invalid user"))
			return
		}

		err = pm.Products.RemoveProduct(c.Request.Context(), &repository.ProductProfileParams{
			ProductID: &req.ProductID,
			StoreId: &requester_id,
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