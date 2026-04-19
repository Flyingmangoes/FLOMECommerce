package models

import "time"

type Store struct {
	StoreId 	string 		`db:"store_id"`
	OwnerId 	string 		`db:"owner_id"`
	StoreName 	string 		`db:"store_name"`
	StoreDesc 	string 		`db:"store_desc"`
	StorePic 	string 		`db:"store_pic"`
	IsActive 	bool 		`db:"is_active"`

	Locale 		*string		`db:"store_locale"`
	Country		*string		`db:"store_country"`
	Address		*string		`db:"store_address"`

	PhoneNumber *string		`db:"store_phone_number"`
	SupportEmail *string	`db:"store_support_email"`

	Instagram	*string 	`db:"store_instagram"`
	Tiktok		*string		`db:"store_tiktok"`
	Website		*string		`db:"store_website"`

	CreatedAt 	time.Time 	`db:"created_at"`
	UpdatedAt 	*time.Time 	`db:"updated_at"`
}