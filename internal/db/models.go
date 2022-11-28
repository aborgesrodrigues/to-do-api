package db

import (
	"database/sql"

	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger
}

type DB struct {
	db     *sql.DB
	logger *zap.Logger
}
