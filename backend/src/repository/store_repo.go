package repository

import (
	"backend/src/models"
	"context"
	"database/sql"
)

type StoreInterface interface {
	CreateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error)

	UpdateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error)
	DeleteStore(ctx context.Context, params *StoreProfileParams) error

	GetStoreByName(ctx context.Context, params *StoreProfileParams) (*models.Store, error)
}

type StoreProfileParams struct {	
	StoreId 		*string 		
	OwnerId 		*string 		
	StoreName 		*string 		
	StoreDesc 		*string 		
	StorePic 		*string 		
	IsActive 		*bool 		
	EmailConsent	*bool
	SmsConsent		*bool
	ConsentSource	*string

	IsAgree			*bool
	IsVerified		*bool	
}

type StoresStore struct {
	db *sql.DB
}



func NewStoresStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}



func (ss *StoresStore) CreateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error) {
	store := &models.Store{}

	tx, err := ss.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_stores (owner_id, store_name, store_desc, store_pic, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING store_id, owner_id, store_name, store_desc, store_pic, is_active, created_at`,
		params.OwnerId, params.StoreName, params.StoreDesc, params.StorePic, params.IsActive,
	).Scan(&store.StoreId, 
		&store.OwnerId, &store.StoreName, &store.StoreDesc, 
		&store.StorePic, &store.IsActive, &store.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return store, nil
}

func (ss *StoresStore) UpdateStore(ctx context.Context, params *StoreProfileParams) (*models.Store, error) {

}

func (ss *StoresStore) DeleteStore(ctx context.Context, params *StoreProfileParams) error {

}

func (ss *StoresStore) GetStoreByName(ctx context.Context, params *StoreProfileParams) (*models.Store, error) {

}