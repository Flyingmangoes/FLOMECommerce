package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
)

type StoreStoreInterface interface {
	CreateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error)
	UpdateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error)
	DeleteStore(ctx context.Context, params *StoreProfileParams) error

	GetStoreByID(ctx context.Context, store_id string) (*models.Store, error)
	FetchStoreID(ctx context.Context, store_id string)(string, error)
	SearchStore(ctx context.Context, params *StoreSearchParams) ([]models.Store, error)
}

type StoreStore struct {
	db *sql.DB
}

func NewStoresStore(db *sql.DB) *StoreStore {
	return &StoreStore{db: db}
}

func (ss *StoreStore) CreateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error) {
	const merchant_type string = "MERCHANT"
	store := &models.Store{}

	tx, err := ss.db.BeginTx(ctx, nil)
	if err != nil { return nil, err }

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_ecommerce.mkt_stores (owner_id, store_name, store_desc, 
		store_pic, store_locale, store_country, 
		store_address, store_phone_number, store_support_email, 
		store_instagram, store_tiktok, store_website)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING store_id, owner_id, store_name, store_desc, store_pic, 
		is_active, store_locale, store_country, store_address, store_phone_number, 
		store_support_email, store_instagram, store_tiktok, store_website, created_at`,
		params.UserId, params.StoreName, params.StoreDesc, 
		params.StorePic, params.Locale, params.Country, params.Address, params.PhoneNumber, params.SupportEmail, params.Instagram, params.Tiktok, params.Website,
	).Scan(
		&store.StoreId, &store.OwnerId, &store.StoreName, &store.StoreDesc, 
		&store.StorePic, &store.IsActive, &store.Locale, 
		&store.Country, &store.Address, &store.PhoneNumber, 
		&store.SupportEmail, &store.Instagram,
		&store.Tiktok, &store.Website, &store.CreatedAt,
	)

	results, err := tx.ExecContext(ctx, 
		`UPDATE mkt_ecommerce.mkt_users SET
			userType 	= COALESCE($1, userType),
			updated_at 	= NOW()
		WHERE user_id = $2`,
		merchant_type, store.OwnerId, 
	)

	rows, err := results.RowsAffected()
	if err != nil { return nil, err }
	if rows == 0 {
		return nil, fmt.Errorf("Not found or Id not match")
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return store, nil
}

func (ss *StoreStore) UpdateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error) {
	store := &models.Store{}

	err := ss.db.QueryRowContext(ctx,
		`UPDATE mkt_ecommerce.mkt_stores SET
			store_name 			= COALESCE ($1, store_name),
			store_desc 			= COALESCE ($2, store_desc),
			store_pic  			= COALESCE ($3, store_pic),
			store_locale 		= COALESCE ($4, store_locale),
			store_country 		= COALESCE ($5, store_country),
			store_address 		= COALESCE ($6, store_address),
			store_phone_number 	= COALESCE ($7, store_phone_number),
			store_support_email = COALESCE ($8, store_support_email),
			store_instagram 	= COALESCE ($9, store_instagram),
			store_tiktok 		= COALESCE ($10, store_tiktok),
			store_website 		= COALESCE ($11, store_website),
			updated_at 			= NOW()
		WHERE store_id = $12 AND owner_id = $13
		RETURNING store_id, owner_id, store_name, store_desc, store_pic, store_locale, 
		store_country, store_address, store_phone_number, store_support_email, store_instagram, 
		store_tiktok, store_website, created_at, updated_at`,
		params.StoreName, params.StoreDesc, params.StorePic, params.Locale, 
		params.Country,params.Address, params.PhoneNumber, 
		params.SupportEmail,params.Instagram, params.Tiktok, 
		params.Website, params.StoreId, params.UserId,
	).Scan(&store.StoreId, &store.OwnerId, 
		&store.StoreName, &store.StoreDesc, &store.StorePic, 
		&store.Locale, &store.Country, &store.Address, 
		&store.PhoneNumber, &store.SupportEmail, &store.Instagram,
		&store.Tiktok, &store.Website, &store.CreatedAt, 
		&store.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (ss *StoreStore) DeleteStore(ctx context.Context, params *StoreProfileParams) error {
	result, err := ss.db.ExecContext(ctx, 
		`DELETE FROM mkt_ecommerce.mkt_stores
		WHERE store_id = $1 AND owner_id = $2`,
		params.StoreId, params.UserId,
	)
	
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("store not found or credentials do not match")
	}

	return nil
}

func (ss *StoreStore) GetStoreByID(ctx context.Context, store_id string) (*models.Store, error) {
	store := &models.Store{}
	var err error
	
	err = ss.db.QueryRowContext(ctx,
		`SELECT store_id, owner_id, store_name, store_desc, store_pic, 
            is_active, store_locale, store_country, store_address, 
            store_phone_number, store_support_email, store_instagram, 
            store_tiktok, store_website, created_at, updated_at
		FROM mkt_ecommerce.mkt_stores WHERE store_id = $1`,
		store_id,
	).Scan(&store.StoreId, &store.OwnerId,
        &store.StoreName, &store.StoreDesc, &store.StorePic,
        &store.IsActive, &store.Locale, &store.Country,
        &store.Address, &store.PhoneNumber, &store.SupportEmail,
        &store.Instagram, &store.Tiktok, &store.Website,
        &store.CreatedAt, &store.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return store, nil
}

func (ss *StoreStore) SearchStore(ctx context.Context, params *StoreSearchParams) ([]models.Store, error) {
	params.Normalize()

	createdAt, id := params.CursorValues()

	rows, err := ss.db.QueryContext(ctx,
		`SELECT store_id, owner_id, store_name, store_desc, store_pic,
				store_country, store_address, store_phone_number,
				created_at, updated_at
		FROM mkt_ecommerce.mkt_stores
		WHERE
			($1::timestamptz IS NULL OR (created_at, store_id) < ($1, $2::uuid)) AND
			($3::varchar IS NULL OR store_name ILIKE '%' || $3 || '%') AND
			($4::varchar IS NULL OR store_country = $4)
		ORDER BY created_at DESC, store_id DESC
		LIMIT $5`,
		createdAt, id,
		params.StoreName, 
		params.StoreCountry, 
		params.Limit + 1,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stores := make([]models.Store, 0)
	for rows.Next(){
		s := models.Store{}
		if err := rows.Scan(
			&s.StoreId, &s.OwnerId, 
			&s.StoreName, &s.StoreDesc,
			&s.StorePic, &s.Country, 
			&s.Address, &s.PhoneNumber,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}

		stores = append(stores, s)
	}

	return stores, rows.Err()
}

func (ss *StoreStore)FetchStoreID(ctx context.Context, store_id string) (string, error) {
	store := models.Store{}

	err := ss.db.QueryRowContext(ctx,
		`SELECT store_id FROM mkt_ecommerce.mkt_stores
		WHERE store_id = $1`,
		store_id,
	).Scan(store.StoreId)

	if err != nil { return "", err }

	return store.StoreId, nil
}
