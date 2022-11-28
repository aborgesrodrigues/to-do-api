package service

import (
	"github.com/aborgesrodrigues/to-do-api/internal/db"
	"go.uber.org/zap"
)

func New(cfg Config) (*Service, error) {
	logger := cfg.Logger

	dbCfg := db.Config{
		Logger: logger,
	}
	db, err := db.New(dbCfg)
	if err != nil {
		logger.Error("Error getting database instance", zap.Error(err))
		return nil, err
	}

	return &Service{
		logger: logger,
		db:     db,
	}, nil
}
