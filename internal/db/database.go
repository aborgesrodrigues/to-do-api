package db

import (
	"database/sql"
	"time"

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

	db.SetMaxOpenConns(10)                  // máximo de conexões abertas
	db.SetMaxIdleConns(5)                   // máximo de conexões ociosas
	db.SetConnMaxLifetime(30 * time.Minute) // tempo máximo de vida da conexão
	db.SetConnMaxIdleTime(5 * time.Minute)  // tempo máximo ocioso

	return &DB{
		logger: logger,
		db:     db,
	}, nil
}
