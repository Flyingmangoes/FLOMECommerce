package models

import "time"

type Order struct {
    ID          string      `db:"order_id"`
    BuyerID     string      `db:"buyer_id"`
    BuyerEmail  string      `db:"buyer_email"`
    TotalPrice  float64     `db:"price_total"`
    Location    string      `db:"location"`
    Status      string      `db:"status"`
    CreatedAt   time.Time   `db:"created_at"`
    OrderItems  []OrderItem
}

type OrderItem struct {
    ID          string  `db:"order_item_id"`
    OrderID     string  `db:"order_id"`
    ProductID   string  `db:"product_id"`
    Quantity    int     `db:"quantity"`
    Price       float64 `db:"price"`     
}