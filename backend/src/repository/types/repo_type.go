package repo_type

import (
	"time"
)

type BaseParams struct {
	UserId 				*string
	Email 				*string
	Username 			*string
	PhoneNumber			*string
}

type UserProfileParams struct {	
	BaseParams
	HashedPassword  	*string
	NewPasswordHashed 	*string
	FirstName			*string
	LastName			*string
	UserType 			*int
	Locale				*string
	Country				*string
	Address 			*string
	EmailConsent		*bool
	SmsConsent			*bool
	ConsentUpdatedAt 	*time.Time
	ConsentSource		*string
	IsAgree				*bool
	IsVerified			*bool	
	CreatedAt			*time.Time
	UpdatedAt			*time.Time
}

type ProductProfileParams struct {
	BaseParams
	ProductID 		*string
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

type OrderParams struct {
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