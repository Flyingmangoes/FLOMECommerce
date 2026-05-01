package services

import (
	"backend/src/repository"
	"context"
	"database/sql"
)

type TxManager struct {
	CentralDB *sql.DB

	Products repository.ProductStoreInterface
	Orders repository.OrderStoreInterface
	Stores repository.StoreStoreInterface
}	

func NewTxManager(db *sql.DB, ps repository.ProductStoreInterface, os repository.OrderStoreInterface, ss repository.StoreStoreInterface ) *TxManager {
	return &TxManager{
		CentralDB: db,
		Products: ps,
		Orders: os,
		Stores: ss,
	}
}

func (tm *TxManager) WithTx(ctx context.Context, fn func (tx *sql.Tx) error) error {
	tx, err := tm.CentralDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}