package services

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ProductStoreInterface interface {
	CreateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
	UpdateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
	RemoveProduct(ctx context.Context, params *ProductProfileParams) error

	GetProductByID(ctx context.Context, lookup *ProductProfileParams) (*models.Product, error)
	GetProductByName(ctx context.Context, lookup *ProductProfileParams) (*models.Product, error)
}

type ProductProfileParams struct {
	//Identifier Section
	ProductID string
	SellerID string
	Storename string
	Name string

	// Profile Section
	Url string
	ImageUrl string
	Price float64
	Rating float64
	Desc string
	Category string
	Availability int

	//Extra Section
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore{
	return &ProductStore{db: db}
}

func (ps *ProductStore) CreateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	newproduct := &models.Product{}
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_products (seller_id, product_name, product_desc, storename, price, product_pic)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING product_id, seller_id, product_name, product_desc, storename, price ,product_pic, created_at`,
	).Scan()

	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newproduct, nil
}

func (ps *ProductStore) UpdateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	updatedprod := &models.Product{}

	err := ps.db.QueryRowContext(ctx,
		`UPDATE mkt_products SET 
			product_name	= COALESCE ($1, product_name),
			product_desc	= COALESCE ($2, product_desc),
			store_name		= COALESCE ($3, store_name),
			product_pic		= COALESCE ($4, product_pic),
			price			= COALESCE ($5, price),
			updated_at		= NOW()
		WHERE product_id = $6
		RETURNING product_name, product_desc, store_name, product_pic, price, updated_at`,
	).Scan()

	if err != nil {
		return nil, err
	}

	return updatedprod, nil
}

func (ps *ProductStore)RemoveProduct(ctx context.Context, params *ProductProfileParams) error {
	results, err := ps.db.ExecContext(ctx,
		`DELETE FROM mkt_products
		WHERE product_id = $1 AND seller_id = $2`,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found or product id not matched")
	}

	return nil
}

func (ps *ProductStore)GetProductByID(ctx context.Context, lookup *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_name, product_desc, seller_id, store_name, product_pic, price, rating, created_at, updated_at FROM mkt_products
		WHERE product_id = $1`,
	).Scan()

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductStore)GetProductByName(ctx context.Context, lookup *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_name, product_desc, seller_id, store_name, product_pic, price, rating, created_at, updated_at FROM mkt_products
		WHERE product_name = $1`,
	).Scan()

	if err != nil {
		return nil, err
	}

	return product, nil 
}