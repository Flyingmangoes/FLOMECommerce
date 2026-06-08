package repository

import (
	"backend/src/models"
	"backend/src/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ProductStoreInterface interface {
	CreateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
	UpdateProduct(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
	RemoveProduct(ctx context.Context, params *ProductProfileParams) error

	GetProductByID(ctx context.Context, product_id string) (*models.Product, error)
	GetProductForUpdate(ctx context.Context, tx *sql.Tx, orderInput []OrderItemInput) ([]models.Product, error)
	FetchStoreID(ctx context.Context, product_id string)(string, error)
	SearchProduct(ctx context.Context, params *ProductSearchParams) ([]models.Product, error)

	UpdateRating(ctx context.Context, params *ProductProfileParams) (*models.Product, error)
    DeductStock(ctx context.Context, tx *sql.Tx, items []OrderItemInput) error
}

type ProductProfileParams struct {
	// Identifier Section
	ProductID 		*string
	StoreId 		*string
	Name 			*string

	// Profile Section
	Url 			*string
	ImageUrl 		*string
	Price 			*float64
	Rating 			*float64
	Desc 			*string
	Category 		*string
	Availability 	*int

	// Extra Section
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ProductSearchParams struct {
	Query 		*string
	Category 	*string
	StoreID 	*string
	MinPrice 	*float64
	MaxPrice 	*float64
	SortBy 		string
	SortOrder 	string
	utils.PagFilter
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
		`INSERT INTO mkt_products (product_name, product_desc, store_id, product_pic, price, category, availability)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING product_id, product_name, product_desc, store_id, product_pic, price, category, rating, availability, created_at`,
		params.Name, params.Desc, params.StoreId, 
		params.ImageUrl, params.Price, params.Category,
		params.Availability,
	).Scan(&product.ProductID, 
		&product.Name, &product.Desc, &product.StoreID, 
		&product.ProductIMG, &product.Price, &product.Category, &product.Rating, 
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
			product_pic		= COALESCE ($3, product_pic),
			price			= COALESCE ($4, price),
			category		= COALESCE ($5, category),
			availability	= COALESCE ($6, availability),
			updated_at		= NOW()
		WHERE product_id = $7 AND store_id = $8
		RETURNING product_id, product_name, product_desc, store_id, product_pic, price, category, availability, updated_at`,
		params.Name, params.Desc, 
	 	params.ImageUrl, params.Price, 
		params.Category, params.Availability, 
		params.ProductID, params.StoreId,
	).Scan(&params.ProductID, &product.Name, 
		&product.Desc, &product.StoreID, &product.ProductIMG, 
		&product.Price, &product.Category, 
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



func (ps *ProductStore)GetProductForUpdate(ctx context.Context, tx *sql.Tx, orderInput []OrderItemInput) ([]models.Product, error) {
	products := []models.Product{}

	for _, p := range orderInput {
		product := &models.Product{}

		err := tx.QueryRowContext(ctx,
			`SELECT product_id, product_name, product_desc, store_id, product_pic, price, rating, availability, created_at, updated_at FROM mkt_products
			WHERE product_id = $1`,
			p.ProductID,
		).Scan(&product.ProductID, &product.Name, 
			&product.Desc, &product.StoreID, 
			&product.ProductIMG, &product.Price, 
			&product.Rating, &product.Availability, &product.CreatedAt, 
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (ps *ProductStore) GetProductByID(ctx context.Context, product_id string) (*models.Product, error) {
	product := &models.Product{}
	
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_id, product_name, product_desc, store_id, product_pic, price, rating, availability, created_at, updated_at 
		FROM mkt_products WHERE product_id = $1`,
		product_id,
	).Scan(&product.ProductID, &product.Name,
		&product.Desc, &product.StoreID,
		&product.ProductIMG, &product.Price,
		&product.Rating, &product.Availability, &product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
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

func (ps *ProductStore) DeductStock(ctx context.Context, tx *sql.Tx, items []OrderItemInput) error {
    for _, item := range items {
        result, err := tx.ExecContext(ctx,
            `UPDATE mkt_products 
            SET availability = availability - $1
            WHERE product_id = $2 AND availability >= $1`,
            item.Quantity, item.ProductID,
        )
		
        if err != nil {
            return err
        }

        rows, err := result.RowsAffected()
        if err != nil {
            return err
        }

        if rows == 0 {
            return fmt.Errorf("insufficient stock for product %s", *item.ProductID)
        }
    }
	
    return nil
}

func (ps *ProductStore) SearchProduct(ctx context.Context, params *ProductSearchParams) ([]models.Product, error) {
	params.Normalize()

	createdAt, id := params.CursorValues()


	rows, err := ps.db.QueryContext(ctx,
		`SELECT product_id, product_name, product_desc, product_pic, store_id, price,
				category, rating, availability, created_at, updated_at
		FROM mkt_ecommerce.mkt_products
		WHERE
			($1::timestamptz IS NULL OR (created_at, product_id) < ($1, $2::uuid))  AND
			($3::varchar IS NULL OR product_name ILIKE '%' || $3 || '%') AND
			($4::varchar IS NULL OR category = $4) 	AND
			($5::uuid IS NULL OR store_id = $5)	 	AND
			($6::numeric IS NULL OR price >= $6) 	AND
			($7::numeric IS NULL OR price <= $7) 	AND
	
			availability > 0
			
		ORDER BY created_at DESC, product_id DESC
		LIMIT $8`,
		createdAt, id,
		params.Query, 
		params.Category,
		params.StoreID,
		params.MinPrice,
		params.MaxPrice,
		params.Limit +1,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next(){
		p := models.Product{}
		if err := rows.Scan(
			p.Name, p.Desc, p.ProductIMG, p.StoreID,
			p.Price, p.Category, p.Rating, p.Availability,
			p.CreatedAt, p.UpdatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, p)
	}
	
	return products, rows.Err()
}

func (ps *ProductStore)FetchStoreID(ctx context.Context, product_id string) (string, error) {
	product := &models.Product{}

	err := ps.db.QueryRowContext(ctx,
		`SELECT store_id FROM mkt_ecommerce.mkt_products 
		WHERE product_id = $1`,
		product_id,
	).Scan(&product.StoreID)

	if err != nil { return "", err}

	return product.StoreID, nil
}