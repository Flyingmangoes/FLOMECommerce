package models

import "time"

type Store struct {
	StoreId 	string 		`db:"store_id"`
	OwnerId 	string 		`db:"owner_id"`
	StoreName 	string 		`db:"store_name"`
	StoreDesc 	string 		`db:"store_desc"`
	StorePic 	string 		`db:"store_pic"`
	IsActive 	bool 		`db:"is_active"`
	CreatedAt 	time.Time 	`db:"created_at"`
	UpdatedAt 	*time.Time 	`db:"updated_at"`
}