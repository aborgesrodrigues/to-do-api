package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBDatabase = "DB_DATABASE"
)

func New(cfg Config) (*DB, error) {
	logger := cfg.Logger
	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv(envDBHost),
		os.Getenv(envDBUser),
		os.Getenv(envDBPassword),
		os.Getenv(envDBDatabase),
		os.Getenv(envDBPort),
	)

	logger.Debug(conn)
	db, err := sql.Open("postgres", conn)
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
