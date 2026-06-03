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
	GetOrders(ctx context.Context, buyer_id, order_id string) ([]*models.Order, []*models.OrderItem, error)
	UpdateOrderStatus(ctx context.Context, order_id, status string) error

	GetOrderId(ctx context.Context, buyer_id string)(string, error)
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
	params.UserId, params.BaseParams.Email, params.TotalPrice, 
	params.Location, params.Status,
	).Scan(&order.ID, 
		&order.BuyerID, &order.BuyerEmail, &order.TotalPrice,
		&order.Location, &order.Status, &order.CreatedAt,
	)

	if err != nil { return nil, nil, err }

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
		if err != nil { return nil, nil, err }
		
		items = append(items, order_item)
	}	

	return order, items, nil
}

func (os *OrderStore)RemoveOrder(ctx context.Context, tx *sql.Tx, params *OrderStoreParams) error {
	results, err := tx.ExecContext(ctx,
		`DELETE FROM mkt_orders
		WHERE order_id = $1 AND buyer_id = $2`,
		params.OrderID, params.UserId,
	)
	if err != nil { return err }

	rows, err := results.RowsAffected()
	if err != nil { return err }
	if rows == 0 {
		return fmt.Errorf("Order not found or id not match")
	}

	return nil
}


func (os *OrderStore) GetOrders(ctx context.Context, buyer_id, order_id string) ([]*models.Order, []*models.OrderItem, error) {
	orders := make([]*models.Order, 0)
	items := make([]*models.OrderItem, 0)

	rows, err := os.db.QueryContext(ctx,
		`SELECT order_id, buyer_id, buyer_email, price_total, location, status, created_at
        FROM mkt_ecommerce.mkt_orders
        WHERE buyer_id = $1 AND order_id = $2
        ORDER BY created_at DESC`,
		buyer_id, order_id,
	)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(
			&order.ID, &order.BuyerID, &order.BuyerEmail,
			&order.TotalPrice, &order.Location, &order.Status,
			&order.CreatedAt,
		); err != nil {
			return nil, nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err();err != nil {
		return nil, nil, err
	}

	for _, order := range orders {
		itemRows, err := os.db.QueryContext(ctx,
			`SELECT order_item_id, order_id, product_id, quantity, price 
			FROM mkt_ecommerce.mkt_order_items WHERE order_id = $1 AND buyer_id = $2`,
			order.ID, order.BuyerID,
		)
		if err != nil {
			return nil, nil, err
		}

		defer itemRows.Close()

		for itemRows.Next(){
			item := &models.OrderItem{}
			if err := itemRows.Scan(
				&item.ID, &item.OrderID, 
				&item.ProductID, &item.Quantity,
				&item.Price,
			); err != nil {
				return nil, nil, err
			}

			items = append(items, item)
		}

		if err := itemRows.Err(); err != nil {
			return nil, nil, err
		}
	}
	
	return orders, items, nil
}

func (os *OrderStore) GetOrderId(ctx context.Context, buyer_id string) (string, error) {
	order := models.Order{}

	err := os.db.QueryRowContext(ctx,
		`SELECT order_id FROM mkt_ecommerce.mkt_orders
		WHERE buyer_id = $1`,
		buyer_id,
	).Scan(order.ID)

	if err != nil {
		return "", err
	}

	return order.ID, nil
}
 
func (os *OrderStore) UpdateOrderStatus(ctx context.Context, order_id, status string) error {
	results, err := os.db.ExecContext(ctx,
		`UPDATE mkt_ecommerce.mkt_orders SET
			status = COALESCE($1, status)
		WHERE order_id = $2 AND buyer_id = $3`,
		status, order_id, 
	)
	if err != nil { return err }

	rows, err := results.RowsAffected()
	if err != nil { return err }
	if rows == 0 {
		return fmt.Errorf("Order not found or id not matched")
	}
	
	return nil
}