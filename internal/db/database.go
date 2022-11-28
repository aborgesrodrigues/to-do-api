package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	envConnString = "CONN_STRING"
)

func New(cfg Config) (*DB, error) {
	logger := cfg.Logger

	logger.Debug(viper.GetString(envConnString))
	db, err := sql.Open("postgres", viper.GetString(envConnString))
	if err != nil {
		logger.Error("Error connecting to database", zap.Error(err))
		return nil, err
	}

	return &DB{
		logger: logger,
		db:     db,
	}, nil
}
