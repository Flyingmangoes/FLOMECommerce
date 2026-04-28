package repository

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

	GetProductByBulk(ctx context.Context, params *ProductProfileParams, limit int) (*models.Product, error)

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
	ProductID *string
	StoreId *string
	Name *string

	// Profile Section
	Url *string
	ImageUrl *string
	Price *float64
	Rating *float64
	Desc *string
	Category *string
	Availability *int

	// Extra Section
	CreatedAt *time.Time
	UpdatedAt *time.Time
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
		`INSERT INTO mkt_products (product_name, product_desc, store_id, url, product_pic, price, category, availability)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING product_id, product_name, product_desc, store_id, url, product_pic, price, category, rating, availability, created_at`,
		params.Name, params.Desc, params.StoreId, params.Url, 
		params.ImageUrl, params.Price, params.Category,
		params.Availability,
	).Scan(&product.ProductID, 
		&product.Name, &product.Desc, &product.StoreID ,&product.Url, 
		&product.ImageUrl, &product.Price, &product.Category, &product.Rating, 
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

	tx, _ := ps.db.BeginTx(ctx, nil)

	defer tx.Rollback()

	err := tx.QueryRowContext(ctx,
		`UPDATE mkt_products SET 
			product_name	= COALESCE ($1, product_name),
			product_desc	= COALESCE ($2, product_desc),
			url				= COALESCE ($3, url),
			product_pic		= COALESCE ($4, product_pic),
			price			= COALESCE ($5, price),
			category		= COALESCE ($6, category),
			availability	= COALESCE ($7, availability),
			updated_at		= NOW()
		WHERE product_id = $8 AND store_id = $9
		RETURNING product_id, product_name, product_desc, store_id, url, product_pic, price, category, availability, updated_at`,
		params.Name, params.Desc, 
		params.Url, params.ImageUrl, params.Price, 
		params.Category, params.Availability, params.ProductID,
		params.StoreId,
	).Scan(&params.ProductID, &product.Name, 
		&product.Desc, &product.StoreID, &product.Url, 
		&product.ImageUrl, &product.Price, &product.Category, 
		&product.Availability, &product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}



func (ps *ProductStore)RemoveProduct(ctx context.Context, params *ProductProfileParams) error {
	results, err := ps.db.ExecContext(ctx,
		`DELETE FROM mkt_products
		WHERE product_id = $1 AND store_id = $2`,
		params.ProductID, params.StoreId,
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
		`SELECT product_id, product_name, product_desc, store_id, url, product_pic, price, rating, availability, created_at, updated_at FROM mkt_products
		WHERE product_id = $1`,
		lookup.ProductID,
	).Scan(&product.ProductID, &product.Name, 
		&product.Desc, &product.StoreID, 
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
		`SELECT product_id, product_name, product_desc, store_id, url, product_pic, price, rating, availability, created_at, updated_at FROM mkt_products
		WHERE product_name = $1`,
		lookup.Name,
	).Scan(&product.ProductID, &product.Name,
		&product.Desc, &product.StoreID,
		&product.Url, &product.ImageUrl, &product.Price,
		&product.Rating, &product.Availability, &product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil 
}



func (ps *ProductStore) GetProductByBulk(ctx context.Context, params *ProductProfileParams, limit int) (*models.Product, error) {
	return nil, nil
}



func (ps *ProductStore) UpdateRating(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}
	
	err := ps.db.QueryRow(
		`UPDATE mkt_products SET
		rating = COALESCE($1, rating)
		WHERE product_id = $2
		RETURNING product_id, rating`,
		params.Rating, params.ProductID,
	).Scan(&product.ProductID, &product.Rating)	

	if err != nil {
		return nil, err
	}

	return product, nil
}



func (ps *ProductStore) UpdateAvailability(ctx context.Context, params *ProductProfileParams) (*models.Product, error) {
	product := &models.Product{}
	
	err := ps.db.QueryRow(
		`UPDATE mkt_products SET
		availability = COALESCE($1, availability)
		WHERE product_id = $2
		RETURNING product_id, availability`,
		params.Availability, params.ProductID,
	).Scan(&product.ProductID, &product.Availability)	

	if err != nil {
		return nil, err
	}

	return product, nil
}

