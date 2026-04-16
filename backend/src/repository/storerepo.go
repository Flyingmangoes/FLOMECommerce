package repository

import (
	"backend/src/models"
	"database/sql"
)

type StoreInterface interface {
	CreateStore() (*models.Store, error)
	UpdateStore() (*models.Store, error)
	DeleteStore() error

	GetStoreByName() (*models.Store, error)
}

type StoreProfileParams struct {	
	
}

type StoresStore struct {
	db *sql.DB
}



func NewStoresStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}