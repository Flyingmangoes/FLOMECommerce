package database

import (
	Logger "backend/src/utils/logger"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDatabaseConnection(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		Logger.Log.Error("detail", zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		Logger.Log.Error("detail", zap.Error(err))
	}

	Logger.Log.Info("Database Established")
	return db
}