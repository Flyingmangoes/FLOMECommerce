package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
)

type CartStoreInterface interface {
	CreateCart(ctx context.Context, params *CartProfileParams)(*models.Cart, error)
	AddItem(ctx context.Context, params *CartProfileParams)(*models.CartItem, error)

	UpdateItemQuantity(ctx context.Context, params *CartProfileParams)(*models.CartItem, error)

	RemoveItem(ctx context.Context, params *CartProfileParams)(error)
	ClearCart(ctx context.Context, params *CartProfileParams)(error)

	GetCart(ctx context.Context, params *CartProfileParams)(*models.Cart, error)
	GetCartItems(ctx context.Context, params []*CartProfileParams)([]*models.CartItem, error)
}

type CartProfileParams struct {
	CartID 		string
	CartItemsID string
	UserID 		string
	ProductID 	string
	Quantity 	int
}

type CartStore struct {
	db *sql.DB
}

func NewCartStore(db *sql.DB) *CartStore{
	return &CartStore{db:db}
}

func (cs *CartStore) CreateCart(ctx context.Context, params *CartProfileParams)(*models.Cart, error) {
	cart := &models.Cart{}

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_ecommerce.mkt_carts (user_id)
		vALUE ($1)
		RETURNING cart_id, user_id`,
		params.UserID,
	).Scan(&cart.ID, &cart.UserID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return cart, nil
}

func (cs *CartStore) UpdateItemQuantity(ctx context.Context, params *CartProfileParams)(*models.CartItem, error){
	updated := &models.CartItem{}

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, 
		`UPDATE mkt_ecommerce.mkt_cart_items SET
			quantity = COALESCE($1, quantity)
		WHERE product_id = $2 AND cart_id = $3
		RETURNING cart_item_id, cart_id, product_id, quantity`,
		params.Quantity, params.ProductID, params.CartID,
	).Scan(&updated.ID, &updated.CartID, &updated.ProductID, &updated.Quantity)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (cs *CartStore) RemoveItem(ctx context.Context, params *CartProfileParams) error {
	results, err := cs.db.ExecContext(ctx,
		`DELETE FROM mkt_ecommerce.mkt_cart_items
		WHERE cart_item_id = $1 AND product_id = $2`,
		params.CartItemsID, params.ProductID,
	) 

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("Item in the card not found or ids provided did not match")
	}
	
	return nil
}

func (cs *CartStore) ClearCart(ctx context.Context, params *CartProfileParams) error {
	results, err := cs.db.ExecContext(ctx,
		`DELETE FROM mkt_ecommerce.mkt_carts
		WHERE cart_id = $1 AND user_id = $2`,
		params.CartID, params.UserID,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("Item in the cart not found or ids provided not match")
	}

	return nil
}

func (cs *CartStore) GetCartItems(ctx context.Context, params []*CartProfileParams)([]*models.CartItem, error) {
	items := make([]*models.CartItem, 0)

	for _, p := range params {
		item := &models.CartItem{}

		err := cs.db.QueryRowContext(ctx,
			`SELECT cart_item_id, cart_id, product_id, quantity
			FROM mkt_ecommerce.mkt_cart_items
			WHERE cart_item_id = $1`,
			p.CartItemsID,
		).Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (cs *CartStore) GetCart(ctx context.Context, params *CartProfileParams) (*models.Cart, error) {
	cart := &models.Cart{}

	err := cs.db.QueryRowContext(ctx,
		`SELECT cart_id, user_id FROM mkt_ecommerce.mkt_carts
		WHERE cart_id = $1`,
		params.CartID,
	).Scan(&cart.ID, &cart.UserID)

	if err != nil {
		return nil, err
	}

	return cart, nil
} 