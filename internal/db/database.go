package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const (
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBDatabase = "DB_DATABASE"
)

var dbExecutionLatencyHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "db_execution_latency",
	Help: "Latency of database execution",
	Buckets: []float64{
		.1,
		1,
		25,
		50,
		100,
		150,
		200,
		250,
		500,
		1000,
		2500,
		5000,
		10000,
	},
}, []string{"table"})

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

	db.SetMaxOpenConns(25)                  // máximo de conexões abertas
	db.SetMaxIdleConns(25)                  // máximo de conexões ociosas
	db.SetConnMaxLifetime(30 * time.Minute) // tempo máximo de vida da conexão
	db.SetConnMaxIdleTime(5 * time.Minute)  // tempo máximo ocioso

	return &DB{
		logger: logger,
		db:     db,
	}, nil
}
