package repository

import "database/sql"

type OrderStoreInterface interface {
	
}

type OrderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) *OrderStore{
	return &OrderStore{db: db}
}

