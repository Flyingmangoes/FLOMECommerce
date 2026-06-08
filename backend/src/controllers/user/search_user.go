package auth_controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchUserRequest struct {
	Query		*string 	`form:"q"`
	Username	*string		`form:"username"`
	SortBy		*string		`form:"sortBy"`
	SortOrder 	*string		`form:"sortOrder"`
	Cursor		*string		`form:"cursor"`
	Limit 		int 		`form:"limit"`
}

func(um *UserManager) SearchUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SearchUserRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Failed to parse client request", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to parse client request"))
			return
		}

		filter := utils.PagFilter{ Limit: req.Limit }
		if req.Cursor != nil {
			cursor, err := utils.DecodeCursor(*req.Cursor)
			if err != nil {
				c.Error(middlewares.ErrUnauthorized("Invalid cursor"))
				return
			}

			filter.Cursor = cursor
		}

		params := &repository.UserSearchParams{
			Query: req.Query,
			Username: req.Username,
			SortBy: req.SortBy,
			SortOrder: req.SortOrder,
			PagFilter: filter,
		}

		filter.Normalize()

		items, err := um.Users.SearchUser(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to retrieve user", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve user"))
			return
		}

		page, err := utils.Build(items, filter.Limit, func(u models.User) (time.Time, string) {
			return u.CreatedAt, u.UserID
		})

		if err != nil {
			Logger.Log.Error("Error while constructing page", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to construct page"))
			return
		}
		c.JSON(http.StatusOK, page)
	}
}