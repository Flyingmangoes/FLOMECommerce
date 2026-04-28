package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
	"time"
)

type OrderStoreInterface interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) (*models.Order, []*models.OrderItem, error)
}

type OrderItemInput struct {
    ProductID string
    Price     float64
    Quantity  int
}

type OrderStoreParams struct {
    OrderID    *string
    BuyerID    *string
    BuyerEmail *string
    TotalPrice *float64
    Status     *string
    Location   *string
	ETA 	   *time.Time
    ProductList []OrderItemInput
}

type OrderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) *OrderStore{
	return &OrderStore{db: db}
}


func (os *OrderStore)CreateOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) (*models.Order, []*models.OrderItem, error) {
	order := &models.Order{}
	items := make([]*models.OrderItem, 0)
	var err error

	err = tx.QueryRowContext(ctx,
	`INSERT INTO mkt_orders(buyer_id, buyer_email, price_total, location, status, eta)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING order_id, buyer_id, buyer_email, price_total, location, status, eta, created_at`,
	params.BuyerID, params.BuyerEmail, params.TotalPrice, 
	params.Location, params.Status, params.ETA,
	).Scan(&order.OrderID, 
		&order.BuyerID, &order.BuyerEmail,&order.TotalPrice,
		&order.Location, &order.Status, &order.ETA,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	for _, item := range params.ProductList {
		order_item := &models.OrderItem{}

    	err = tx.QueryRowContext(ctx,
        	`INSERT INTO mkt_orders_item(order_id, product_id, quantity, price)
        	VALUES($1, $2, $3, $4)
        	RETURNING order_item_id, order_id, product_id, quantity, price`,
        	order.OrderID, item.ProductID, item.Quantity, item.Price,
    	).Scan(
			&order_item.OrderItemID, &order_item.OrderID, &order_item.ProductID,
			&order_item.Quantity, &order_item.Price,
		)

		if err != nil {
			return nil, nil, err
		}

		items = append(items, order_item)
	}	

	return order, items, nil
}