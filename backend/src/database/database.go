package database

import (
	logger_system "backend/src/utils/LoggerSystem"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDatabaseConnection(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger_system.Log.Error("detail", zap.Error(err))
	}

	if err = db.Ping(); err != nil {
		logger_system.Log.Error("detail", zap.Error(err))
	}

	logger_system.Log.Info("Database Established")
	return db
}