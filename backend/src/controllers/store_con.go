package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	repo "backend/src/repository"
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

/* STORE DOCUMENTATION
*	StoreManager is just a struct for all the required repos for store to be functional
*	it used in backend/src/server/routing.go
*
*	I MOVE EVERY REQUEST STRUCT FROM ALL CONTROLLER TO requers_variant.go
*
*	To make a user became a merchant this is how to do it:
*		1. Register new account / or login to an account
*		2. Send the required data to RegisterStore
*
*	1. Register Store
*		This function handle merchant registry, it require a Request struct.
*		This struct hold the json data required to register a new merchant account.
*
*		The request then be bind to json it return an error or nil if succeed.
*		Next we filled up StoreProfileParams(required for CreateUser repository)
*		with the data in request.
*
*		If the CreateUser succeed it will return *models.Store and nil, if failed
*		it will return nil and err
*
*		After that just send the result using c.JSON(http.StatusCreated, result)
*
*	2. Update Store
*		
*
*
*
*
*
*	3. Remove Store
*
*
*
*
*
*
 */

type StoreManager struct {
	Users 		repo.UserStoreInterface
	Stores    	repo.StoreStoreInterface
    Tokens   	repo.TokenStoreInterface
	Cache 		redis.RedisInterface
	JWTSecret 	[]byte
}

func (sm *StoreManager) RegisterStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StoreRegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request:"))
			return
		}

		ownerid := c.GetString("userId")

		params := &repo.StoreProfileParams{
			BaseParams: repo.BaseParams{
				UserId: &ownerid,
				Locale: &req.Locale,
				Country: &req.Country,
				Address: &req.Address,
			},

			StoreName: &req.StoreName,
			StoreDesc: &req.StoreDesc,
			StorePic: &req.StoreIMG,
			PhoneNumber: &req.PhoneNumber,
			SupportEmail: &req.SupportEmail,
			Instagram: &req.Instagram,
			Tiktok: &req.Tiktok,
			Website: &req.Website,
		}

		store, err := sm.Stores.CreateStore(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to create stores"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "store created",
				"store": store,
			},
		})
	}
}

func (sm *StoreManager) UpdateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StoreUpdateRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		ownerId := c.GetString("userId")
        storeId := c.GetString("storeId") 

		hashed_password, err := sm.Users.GetPassword(c.Request.Context(), &repo.UserProfileParams{
			BaseParams: repo.BaseParams{ UserId: &ownerId },
		})
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to compare credentials"))
			return
		}

		if err = validators.ValidatePassword(*hashed_password, req.OwnerPassword); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		params:= &repo.StoreProfileParams{	
			BaseParams: repo.BaseParams{
				Locale: req.NewLocale,
				Country: req.NewCountry,
				Address: req.NewAddress,
			},

			StoreId: &storeId,
			StoreName: req.NewName,
			StoreDesc: req.NewDesc,
			StorePic:  req.NewImage,
			PhoneNumber: req.NewPhoneNumber,
			SupportEmail: req.NewSupportEmail,
			Instagram: req.NewInstagram,
			Tiktok: req.NewTiktok,
			Website: req.NewWebsite,
		}

		store, err := sm.Stores.UpdateStore(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update store"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "store updated",
				"store": store,
			},
		})
	}	
}

func (sm *StoreManager) DeleteStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StoreRemoveRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("Error in reading client request", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to parse client request"))
			return
		}

		requester_id := c.GetString("userId")
		isAllowed, err := validators.ValidateRequester(c.Request.Context(), 
			requester_id, req.StoreId, "PRODUCT", sm.Stores,
		)
		if err != nil {
			Logger.Log.Error("Failed in Validate request", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to compare credentials"))
			return
		}

		if !isAllowed {
			c.Error(middlewares.ErrUnauthorized("Invalid user"))
			return
		}

		params := &repo.StoreProfileParams{
			StoreId: &req.StoreId,
			BaseParams: repo.BaseParams{
				UserId: &requester_id,
			},
		}

		err = sm.Stores.DeleteStore(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Error in deletion process", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to process deletion"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": "store removed",
		})
	}
}

func (sm *StoreManager) SearchStore() gin.HandlerFunc{
	return func (c *gin.Context) {
		var req SearchStoreRequest

		if err := c.ShouldBindQuery(&req);err != nil {
			Logger.Log.Error("Failed to read request", zap.Error(err))
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

		params := &repo.StoreSearchParams{
			Query: req.Query,
			StoreName: req.StoreName,
			StoreCountry: req.StoreCountry,
			SortBy: req.SortBy,
			SortOrder: req.SortOrder,
			PagFilter: filter,
		}

		url, _ := url.Parse(fmt.Sprintf("stores:%s", c.Request.URL))
		cacheKey := sm.Cache.GenerateCacheKey(url)

		cached, err := sm.Cache.Get(c.Request.Context(), cacheKey)
		if err == nil && cached != nil {
			var page utils.Page[models.Store]
			if err := json.Unmarshal(cached, &page); err != nil {
				c.JSON(http.StatusOK, page)
				return
			}
		} 

		stores, err := sm.Stores.SearchStore(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Failed to retrieve store", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to retrieve Stores"))
			return
		}

		page, err := utils.Build(stores, filter.Limit, func(s models.Store) (time.Time, string)  {
			return s.CreatedAt, s.StoreId
		})
		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to build page"))
			return
		}
		c.JSON(http.StatusOK, page)
	}
}