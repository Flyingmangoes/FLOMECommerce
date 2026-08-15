package repository

import (
	"backend/src/utils"
	"time"
)

type BaseParams struct {
	UserId 				*string
	Email 				*string
	Username 			*string
	Locale				*string
	Country				*string
	Address 			*string
}

type UserProfileParams struct {	
	BaseParams
	PhoneNumber			*string
	HashedPassword  	*string
	NewPasswordHashed 	*string
	FirstName			*string
	LastName			*string
	UserType 			*int
	EmailConsent		*bool
	SmsConsent			*bool
	ConsentUpdatedAt 	*time.Time
	ConsentSource		*string
	IsAgree				*bool
	IsVerified			*bool	
	CreatedAt			*time.Time
	UpdatedAt			*time.Time
}

type UserSearchParams struct {
	Query 				*string
	Username 			*string
	SortBy 				*string
	SortOrder 			*string
}

type StoreProfileParams struct {
	BaseParams
    StoreId      *string
    StoreName    *string
    StoreDesc    *string
    StorePic     *string
    IsActive     *bool
    PhoneNumber  *string
    SupportEmail *string
    Instagram    *string
    Facebook     *string
    Tiktok       *string
    Website      *string
}

type StoreSearchParams struct {
	Query 			*string	
	StoreName 		*string 	
	StoreCountry 	*string 	
	SortBy 			*string		
	SortOrder 		*string		
}

type ProductProfileParams struct {
	ProductID 		*string
	StoreId 		*string
	Name 			*string
	ImageUrl 		*string
	Price 			*float64
	Rating 			*float64
	Desc 			*string
	Category 		*string
	Availability 	*int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ProductSearchParams struct {
	Query 		*string
	Category 	*string
	StoreID 	*string
	MinPrice 	*float64
	MaxPrice 	*float64
	SortBy 		*string
	SortOrder 	*string
}

type OrderItemInput struct {
    ProductID *string
    Price     *float64
    Quantity  *int
}

type OrderStoreParams struct {
	BaseParams
    OrderID    *string
    TotalPrice *float64
    Status     *string
    Location   *string
    ProductList []OrderItemInput
}

type CartProfileParams struct {
	BaseParams
	CartID 		*string
	CartItemsID *string
	ProductID 	*string
	StoreID 	*string
	Quantity 	*int
}