package product

import (
	product_types "backend/src/controllers/product/types"
	terror "backend/src/error"
	"backend/src/models"
	repo_type "backend/src/repository/types"
	logger_system "backend/src/utils/LoggerSystem"
	pagination "backend/src/utils/Pagination"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (pm *ProductManager) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req product_types.SearchProductRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			logger_system.Log.Error("Failed to bind query", zap.Error(err))
			c.Error(terror.ErrBadRequest("Failed to read client request"))
			return
		}

		filter := pagination.PagFilter{Limit: req.Limit}

		if req.Cursor != nil {
			cursor, err := pagination.DecodeCursor(*req.Cursor)
			if err != nil {
				c.Error(terror.ErrBadRequest("Invalid cursor"))
				return 
			}

			filter.Cursor = cursor
		}

		filter.Normalize()

		params := &repo_type.ProductSearchParams{
			Query: req.Query,
			Category: req.Category,
			StoreID: req.StoreID,
			MinPrice: req.MinPrice,			
			MaxPrice: req.MaxPrice,
			SortBy: req.SortBy,
			SortOrder: req.SortOrder,
		}

		cacheKey := pm.Cache.GenerateCacheKey("products", c.Request.URL.String())
		cached, err := pm.Cache.Get(c.Request.Context(), cacheKey)
		logger_system.Log.Info("cached", zap.Any("cached value", cached))
		if err == nil && cached != nil {
			var page pagination.Page[models.Product]
			if err :=  json.Unmarshal(cached, &page); err == nil {
				logger_system.Log.Info("build", zap.Any("page", page))
				c.JSON(http.StatusOK, page)
				return
			}
		}

		products, err := pm.Products.Search(c.Request.Context(), params)
		if err != nil {
			logger_system.Log.Error("Failed to retrieve products", zap.Error(err))
			c.Error(terror.ErrInternal("Failed to retrieve products"))
			return
		}

		page, err := pagination.Build(products, filter.Limit, func(p models.Product) (time.Time, string) {
			return p.CreatedAt, p.ProductID			
		})
		if err != nil {
			c.Error(terror.ErrInternal("Failed to build page"))
			return
		}

		logger_system.Log.Info("products", zap.Any("array", products))
		logger_system.Log.Info("build", zap.Any("page", page))

		cacheValue, err := json.Marshal(page)
		if err != nil {
			logger_system.Log.Debug("Failed to marshal page", zap.Error(err))
		} else {
			if err := pm.Cache.Set(c.Request.Context(), cacheKey, string(cacheValue)); err != nil {
				logger_system.Log.Debug("Failed to set cache", zap.Error(err))
			}
		}

		c.JSON(http.StatusOK, page)
	}
}