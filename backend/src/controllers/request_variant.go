package controllers

//
// MERCHANT REQUEST
//

type StoreRegisterRequest struct {	
	StoreName 	 string `json:"storeName"`
	StoreDesc 	 string `json:"storeDesc"`
	StoreIMG 	 string `json:"storeIMG"`

	Locale 		 string	`json:"storeLocale"`
	Country		 string	`json:"storeCountry"`
	Address		 string	`json:"storeAddress"`

	PhoneNumber  string	`json:"storePhoneNumber"`
	SupportEmail string	`json:"storeSupportEmail"`

	Instagram	 string `json:"storeInstagram"`
	Tiktok		 string	`json:"storeTiktok"`
	Website		 string	`json:"storeWebsite"`	
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
	StoreId		string 	`json:"storeId"`
}

type SearchStoreRequest struct {
	Query 			*string		`form:"q"`
	StoreName		*string		`form:"storeName"`
	StoreCountry	*string		`form:"country"`
	SortBy			*string		`form:"sortBy"`
	SortOrder 		*string		`form:"sortOrder"`
	Cursor			*string		`form:"cursor"`
	Limit 			int			`form:"limit"`
}

//
// PRODUCT REQUEST
//

type RegisterProductRequest struct {
	ProductName 	string 	`json:"productName"`
	ImageUrl 		string 	`json:"imageUrl"`
	Price 			float64 `json:"price"`
	Desc 			string  `json:"desc"`
	Category 		string	`json:"category"`
	Availability 	int		`json:"availability"`
}

type UpdateProductRequest struct {
	ProductID		string 		`json:"productId" binding:"required"`
	NewProductName 	*string 	`json:"newProductName" binding:"omitempty"`
	NewProductDesc 	*string		`json:"newProductDesc" binding:"omitempty"`
	NewStorename 	*string		`json:"newStoreName" binding:"omitempty"`
	NewImage 		*string 	`json:"newImage" binding:"omitempty"`
	NewPrice 		*float64 	`json:"newPrice" binding:"omitempty"`
	NewCategory 	*string		`json:"newCategory" binding:"omitempty"`
	NewAvailability *int		`json:"newAvailability" binding:"omitempty"`
}

type RemoveProductRequest struct {
	ProductID 		string	`json:"productId" binding:"required"`
}

type SearchProductRequest struct {
	Query 		*string 	`form:"q"`
	Category 	*string 	`form:"category"`
	StoreID		*string 	`form:"storeId"`
	MinPrice	*float64 	`form:"minPrice"`
	MaxPrice 	*float64 	`form:"maxPrice"`
	SortBy		*string		`form:"sortBy"`
	SortOrder 	*string		`form:"sortOrder"`
	Cursor		*string  	`form:"cursor"`
	Limit 		int			`form:"limit"`
}

//
//	ORDER REQUEST
//

type OrderItemRequest struct {
    ProductID 	string `json:"productId" binding:"required"`
    Quantity  	int    `json:"quantity"  binding:"required,min=1"`
}

type OrderRequest struct {
	Items []OrderItemRequest 	`json:"items" binding:"required"`
    BuyerLocation *string       `json:"location" binding:"omitempty"`
}	

type CancelOrderRequest struct {
	OID string `json:"orderId" binding:"required"`
}

//
//	CART REQUEST
//

type AddItemRequest struct {
	ProductID 	string 	`json:"productId" binding:"required"`
	Quantity  	int 	`json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	CartItemId 	string `json:"cartItemId" binding:"required"`
}

type UpdateQuantityRequest struct {
	CartItemID 	string	`json:"cartItemId" binding:"required"`
	NewQuantity int 	`json:"newQuantity" binding:"required"`
}