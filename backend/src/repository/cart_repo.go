package repository

import (
	"backend/src/models"
	repo_type "backend/src/repository/types"
	"context"
	"database/sql"
)

type CartStoreInterface interface {
	Create(ctx context.Context, user_id string)(*models.Cart, error)
	InsertItems(ctx context.Context, params *repo_type.CartProfileParams)(*models.CartItem, error)
	UpdateQuantity(ctx context.Context, params *repo_type.CartProfileParams)(*models.CartItem, error)

	RemoveItems(ctx context.Context, params *repo_type.CartProfileParams)(error)
	ClearItems(ctx context.Context, params *repo_type.CartProfileParams)(error)
	Get(ctx context.Context, params *repo_type.CartProfileParams)(*models.Cart, error)
	GetItems(ctx context.Context, params *repo_type.CartProfileParams)([]*models.CartItem, error)
}

type CartStore struct {
	db *sql.DB
}

func NewCartStore(db *sql.DB) *CartStore{
	return &CartStore{db:db}
}

func (cs *CartStore) Create(ctx context.Context, user_id string)(*models.Cart, error) {
	cart := &models.Cart{}

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_ecommerce.mkt_carts (user_id)
		VALUES ($1)
		RETURNING cart_id, user_id
		ON CONFLICT(user_id)
		DO RETURNING cart_id, user_id`,
		user_id,
	).Scan(&cart.ID, &cart.UserID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return cart, nil
}

func (cs *CartStore) InsertItems(ctx context.Context, params *repo_type.CartProfileParams)(*models.CartItem, error) {
	cart_item := &models.CartItem{}

	err := cs.db.QueryRowContext(ctx,
		`INSERT INTO mkt_ecommerce.mkt_cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT(cart_id, product_id)
		DO UPDATE SET 
			quantity = mkt_ecommerce.mkt_cart_items.quantity + EXCLUDED.quantity
		RETURNING cart_item_id, cart_id, product_id, quantity`,
		params.CartID, params.ProductID, params.Quantity,
	).Scan(
		&cart_item.ID, &cart_item.CartID, 
		&cart_item.ProductID, &cart_item.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return cart_item, nil
}

func (cs *CartStore) UpdateQuantity(ctx context.Context, params *repo_type.CartProfileParams)(*models.CartItem, error){
	updated := &models.CartItem{}

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, 
		`UPDATE mkt_ecommerce.mkt_cart_items SET
			quantity = COALESCE($1, quantity)
		WHERE cart_item_id = $2 AND cart_id = $3
		RETURNING cart_item_id, cart_id, product_id, quantity`,
		params.Quantity, params.CartItemsID, params.CartID,
	).Scan(&updated.ID, &updated.CartID, 
		&updated.ProductID, &updated.Quantity,
	)
	
	if err != nil { return nil, err }
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (cs *CartStore) RemoveItems(ctx context.Context, params *repo_type.CartProfileParams) error {
	results, err := cs.db.ExecContext(ctx,
		`DELETE FROM mkt_ecommerce.mkt_cart_items
		WHERE cart_item_id = $1 AND cart_id = $2`,
		params.CartItemsID, params.CartID,
	)     

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return nil
	}
	
	return nil
}

func (cs *CartStore) ClearItems(ctx context.Context, params *repo_type.CartProfileParams) error {
	results, err := cs.db.ExecContext(ctx,
		`DELETE FROM mkt_ecommerce.mkt_cart_items
		WHERE cart_id = $1`,
		params.CartID,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return nil
	}

	return nil
}

func (cs *CartStore) GetItems(ctx context.Context, params *repo_type.CartProfileParams)([]*models.CartItem, error) {
	items := make([]*models.CartItem, 0)

	rows, err := cs.db.QueryContext(ctx,
		`SELECT cart_item_id, cart_id, product_id, quantity
		FROM mkt_ecommerce.mkt_cart_items
		WHERE cart_id = $1`,
		params.CartID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.CartItem{}
		if err := rows.Scan(
				&item.ID, &item.CartID, 
				&item.ProductID, &item.Quantity,
			); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (cs *CartStore) Get(ctx context.Context, params *repo_type.CartProfileParams) (*models.Cart, error) {
	cart := &models.Cart{}

	err := cs.db.QueryRowContext(ctx,
		`SELECT cart_id, user_id FROM mkt_ecommerce.mkt_carts
		WHERE user_id = $1`,
		params.UserId,
	).Scan(&cart.ID, &cart.UserID)

	if err != nil {
		return nil, err
	}

	return cart, nil
} 