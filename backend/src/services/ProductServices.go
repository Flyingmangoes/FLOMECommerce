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

	GetProductList(ctx context.Context, params *ProductProfileParams, limit int) (*models.Product, error)

	UpdateRating(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
	UpdateAvailability(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
}

type ProductProfileParams struct {
	// Use the same structure for updating a product,
	// product id cannot be updated, it stay forever like 
	// how it supposed to be, when updating product id used 
	// for specified which product to update
	//
	// keep that in mind for whoever find this usefull

	// Identifier Section
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

	// Extra Section
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
	product := &models.Product{}
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_products (product_id, seller_id, product_name, product_desc, storename, url, product_pic, price, category, availability)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING product_id, seller_id, product_name, product_desc, storename, url, product_pic, price, category, rating, availability, created_at, updated_at`,
		params.ProductID, params.SellerID, params.Name, 
		params.Desc, params.Storename, params.Url, 
		params.ImageUrl, params.Price, params.Category,
		params.Availability,
	).Scan(&product.ProductID, 
		&product.SellerID, &product.StoreName, 
		&product.Name, &product.Desc, &product.Url, 
		&product.ImageUrl, &product.Category, &product.Rating, 
		&product.Availability, &product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}



func (ps *ProductStore) UpdateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}

	err := ps.db.QueryRowContext(ctx,
		`UPDATE mkt_products SET 
			product_name	= COALESCE ($1, product_name),
			product_desc	= COALESCE ($2, product_desc),
			store_name		= COALESCE ($3, store_name),
			url				= COALESCE ($4, url),
			product_pic		= COALESCE ($5, product_pic),
			price			= COALESCE ($6, price),
			category		= COALESCE ($7, category),
			availability	= COALESCE ($8, availability),
			updated_at		= NOW()
		WHERE product_id = $9
		RETURNING product_name, product_desc, store_name, url, product_pic, price, category, availability, updated_at`,
		params.Name, params.Desc, params.Storename, 
		params.Url, params.ImageUrl, params.Price, 
		params.Category, params.Availability, params.ProductID,
	).Scan(&product.Name, 
		&product.Desc, &product.StoreName, &product.Url, 
		&product.ImageUrl, &product.Price, &product.Category, 
		&product.Availability, &product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}



func (ps *ProductStore)RemoveProduct(ctx context.Context, params *ProductProfileParams) error {
	results, err := ps.db.ExecContext(ctx,
		`DELETE FROM mkt_products
		WHERE product_id = $1 AND seller_id = $2`,
		params.ProductID, params.SellerID,
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
		`SELECT ,product_id, product_name, product_desc, seller_id, store_name, url, product_pic, price, rating, availability, created_at, updated_at FROM mkt_products
		WHERE product_id = $1`,
		lookup.ProductID,
	).Scan(&product.ProductID, &product.Name, 
		&product.Desc, &product.SellerID, &product.StoreName, 
		&product.Url, &product.ImageUrl, &product.Price, 
		&product.Rating, &product.Availability, &product.CreatedAt, 
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}



func (ps *ProductStore)GetProductByName(ctx context.Context, lookup *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}
	
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_id, product_name, product_desc, seller_id, store_name, url, product_pic, price, rating, availability, created_at, updated_at FROM mkt_products
		WHERE product_name = $1`,
		lookup.Name,
	).Scan(&product.ProductID, &product.Name,
		&product.Desc, &product.SellerID, &product.StoreName,
		&product.Url, &product.ImageUrl, &product.Price,
		&product.Rating, &product.Availability, &product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil 
}



func (ps *ProductStore) GetProductList(ctx context.Context, params *ProductProfileParams, limit int) (*models.Product, error) {
	return nil, nil
}



func (ps *ProductStore) UpdateRating(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	return nil, nil
}



func (ps *ProductStore) UpdateAvailability(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	return nil, nil
}


