package controllers

import (
	"backend/src/middlewares"
	repo "backend/src/repository"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"backend/src/validators"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//
// STORE HANDLER DECLARATION
//

type StoreManager struct {
	Users 		repo.UserStoreInterface
	Stores    	repo.StoreStoreInterface
    Tokens   	repo.TokenStoreInterface
	JWTSecret 	[]byte
}

type StoreRegisterRequest struct {	
	StoreName 	 string `json:"storeName"`
	StoreDesc 	 string `json:"storeDesc"`
	StoreIMG 	 string `json:"storeIMG"`

	Locale 		 string	`json:"storeLocale"`
	Country		 string	`json:"storeCountry"`
	Address		 string	`json:"storeAddress"`

	PhoneNumber  string	`json:"storePhoneNumber"`
	SupportEmail string	`json:"storeSupportEmail"`

	Instagram	string 	`json:"storeInstagram"`
	Tiktok		string	`json:"storeTiktok"`
	Website		string	`json:"storeWebsite"`	
}

type StoreUpdateRequest struct {
	OwnerPassword 	string `json:"password" binding:"required"`

	NewName    		*string `json:"newStoreName" binding:"omitempty"`
	NewDesc 		*string `json:"newStoreDesc" binding:"omitempty"`
	NewImage		*string `json:"newStorePic" binding:"omitempty"`
	NewPhoneNumber  *string `json:"newPhoneNumber" binding:"omitempty"`
	NewSupportEmail *string `json:"newSupportEmail" binding:"omitempty"`

    NewLocale  		*string `json:"newLocale" binding:"omitempty"`
	NewCountry		*string `json:"newCountry" binding:"omitempty"`
	NewAddress		*string `json:"newAddress" binding:"omitempty"`

	NewTiktok    	*string `json:"newTiktokAcc" binding:"omitempty"`
	NewInstagram 	*string `json:"newInstagramAcc" binding:"omitempty"`
	NewWebsite 		*string `json:"newWebsite" binding:"omitempty"`
}

type StoreRemoveRequest struct {
	OwnerId		string  `json:"userId"      binding:"required"`
	Password	string  `json:"password"     binding:"required"`
}

//
// STORE HANDLER IMPLEMENTATION
//

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

		c.JSON(http.StatusOK, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": "store removed",
		})
	}
}



