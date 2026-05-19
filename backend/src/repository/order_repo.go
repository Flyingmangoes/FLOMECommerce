package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
)

type OrderStoreInterface interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) (*models.Order, []*models.OrderItem, error)
	RemoveOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) error
}

type OrderItemInput struct {
    ProductID *string
    Price     *float64
    Quantity  *int
}

type OrderStoreParams struct {
    OrderID    *string
    BuyerID    *string
    BuyerEmail *string
    TotalPrice *float64
    Status     *string
    Location   *string
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
	`INSERT INTO mkt_orders(buyer_id, buyer_email, price_total, location, status)
	VALUES($1, $2, $3, $4, $5)
	RETURNING order_id, buyer_id, buyer_email, price_total, location, status, created_at`,
	params.BuyerID, params.BuyerEmail, params.TotalPrice, 
	params.Location, params.Status,
	).Scan(&order.ID, 
		&order.BuyerID, &order.BuyerEmail, &order.TotalPrice,
		&order.Location, &order.Status, &order.CreatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	for _, item := range params.ProductList {
		order_item := &models.OrderItem{}

    	err = tx.QueryRowContext(ctx,
        	`INSERT INTO mkt_order_items(order_id, product_id, quantity, price)
        	VALUES($1, $2, $3, $4)
        	RETURNING order_item_id, order_id, product_id, quantity, price`,
        	order.ID, item.ProductID, item.Quantity, item.Price,
    	).Scan(
			&order_item.ID, &order_item.OrderID, &order_item.ProductID,
			&order_item.Quantity, &order_item.Price,
		)

		if err != nil {
			return nil, nil, err
		}

		items = append(items, order_item)
	}	

	return order, items, nil
}

func (os *OrderStore)RemoveOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) error {
	results, err := tx.ExecContext(ctx,
		`DELETE FROM mkt_orders
		WHERE order_id = $1 AND buyer_id = $2`,
		params.OrderID, params.BuyerID,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return fmt.Errorf("orders not found or orders id is not match")
	}

	return nil
}