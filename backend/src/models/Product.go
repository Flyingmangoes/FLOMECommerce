package models

import (
    "time"
)

type Product struct {
    ProductID       string      `db:"product_id"`
    Name            string      `db:"product_name"`
    Desc            string      `db:"product_desc"`
    ProductIMG      string      `db:"product_pic"`  
    Price           float64     `db:"price"`
    Category        string      `db:"category"`
    Rating          float64     `db:"rating"`
    Availability    int         `db:"availability"`
    CreatedAt       time.Time   `db:"created_at"` 
    UpdatedAt       *time.Time  `db:"updated_at"`
}