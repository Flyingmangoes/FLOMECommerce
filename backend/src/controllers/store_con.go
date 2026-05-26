package controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
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
	Users 		repository.UserStoreInterface
	Stores    	repository.StoreStoreInterface
    Tokens   	repository.TokenStoreInterface
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
	Confirmation 	bool   `json:"confirmation" binding:"required"`

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

		params := &repository.StoreProfileParams{
			OwnerId: &ownerid,
			StoreName: &req.StoreName,
			StoreDesc: &req.StoreDesc,
			StorePic: &req.StoreIMG,
			Locale: &req.Locale,
			Country: &req.Country,
			Address: &req.Address,
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

		c.JSON(http.StatusCreated, gin.H{"response": store})
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

		up := &repository.UserProfileParams{
			UserId: &ownerId,
		}

		userpass, err := sm.Users.GetPassword(c.Request.Context(), up)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to compare credentials"))
			return
		}

		if err = validators.ValidatePassword(userpass.PasswordHash, req.OwnerPassword); err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		params:= &repository.StoreProfileParams{	
			StoreId: &storeId,
			StoreName: req.NewName,
			StoreDesc: req.NewDesc,
			StorePic:  req.NewImage,
			Locale: req.NewLocale,
			Country: req.NewCountry,
			Address: req.NewAddress,
			PhoneNumber: req.NewPhoneNumber,
			SupportEmail: req.NewSupportEmail,
			Instagram: req.NewInstagram,
			Tiktok: req.NewTiktok,
			Website: req.NewWebsite,
			Confirmation: req.Confirmation,
		}

		store, err := sm.Stores.UpdateStore(c.Request.Context(), params)
		if err != nil {
			Logger.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update store"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"response": store})
	}	
}



func (sm *StoreManager) DeleteStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}



