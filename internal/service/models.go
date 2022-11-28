package service

import (
	"github.com/aborgesrodrigues/to-do-api/internal/db"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger
}

type Service struct {
	Logger *zap.Logger
	DB     *db.DB
}
