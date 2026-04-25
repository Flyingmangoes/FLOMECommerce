package controllers

import (
	"backend/src/middlewares"
	"backend/src/repository"
	"backend/src/utils"
	"backend/src/validators"

	"github.com/gin-gonic/gin"
	"net/http"
	"go.uber.org/zap"
)

type StoreContext struct {
	Users 		repository.UserStoreInterface
	Stores    	repository.StoreStoreInterface
    Tokens   	repository.TokenStoreInterface
	JWTSecret 	[]byte
}

type StoreRegisterRequest struct {	
	StoreName 	 string `json:"storeName"`
	StoreDesc 	 string `json:"storeDesc"`
	StorePic 	 string `json:"storePic"`

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

	NewStoreName    *string `json:"newStoreName" binding:"omitempty"`
	NewStoreDesc 	*string `json:"newStoreDesc" binding:"omitempty"`
	NewStorePic		*string `json:"newStorePic" binding:"omitempty"`
	NewPhoneNumber  *string `json:"newPhoneNumber" binding:"omitempty"`
	NewSupportEmail *string `json:"newSupportEmail" binding:"omitempty, email"`

    NewLocale  		*string `json:"newLocale" binding:"omitempty"`
	NewCountry		*string `json:"newCountry" binding:"omitempty"`
	NewAddress		*string `json:"newAddress" binding:"omitempty"`

	NewTiktokAcc    *string `json:"newTiktokAcc" binding:"omitempty"`
	NewInstagramAcc *string `json:"newInstagramAcc" binding:"omitempty"`
	NewWebsite 		*string `json:"newWebsite" binding:"omitempty"`
}

type StoreLoginRequest struct {
	OwnerEmail string `json:"ownerEmail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type StoreRemoveRequest struct {
	OwnerId		string  `json:"userId"      binding:"required"`
	Password	string  `json:"password"     binding:"required"`
}

func (sc *StoreContext) RegisterStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StoreRegisterRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request:"))
			return
		}

		ownerid := c.GetString("userId")

		params := &repository.StoreProfileParams{
			OwnerId: &ownerid,
			StoreName: &req.StoreName,
			StoreDesc: &req.StoreDesc,
			StorePic: &req.StorePic,
			Locale: &req.Locale,
			Country: &req.Country,
			Address: &req.Address,
			PhoneNumber: &req.PhoneNumber,
			SupportEmail: &req.SupportEmail,
			Instagram: &req.Instagram,
			Tiktok: &req.Tiktok,
			Website: &req.Website,
		}

		store, err := sc.Stores.CreateStore(c.Request.Context(), params)
		if err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to create stores"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"store": store})
	}
}

func (sc *StoreContext) UpdateStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StoreUpdateRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		ownerid := c.GetString("userId")
        storeid := c.GetString("storeId") 

		up := &repository.UserProfileParams{
			UserId: &ownerid,
		}

		userpass, err := sc.Users.GetPassword(c.Request.Context(), up)
		if err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to compare credentials"))
			return
		}

		if err = validators.ValidatePassword(userpass.PasswordHash, req.OwnerPassword); err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		params:= &repository.StoreProfileParams{	
			StoreName: req.NewStoreName,
			StoreDesc: req.NewStoreDesc,
			StorePic:  req.NewStorePic,
			Locale: req.NewLocale,
			Country: req.NewCountry,
			Address: req.NewAddress,
			PhoneNumber: req.NewPhoneNumber,
			SupportEmail: req.NewSupportEmail,
			Instagram: req.NewInstagramAcc,
			Tiktok: req.NewTiktokAcc,
			Website: req.NewWebsite,
			StoreId: &storeid,
		}

		store, err := sc.Stores.UpdateStore(c.Request.Context(), params)
		if err != nil {
			utils.Log.Error("Error", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to update store"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"store": store})
	}	
}

func (sc *StoreContext) DeleteStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}



