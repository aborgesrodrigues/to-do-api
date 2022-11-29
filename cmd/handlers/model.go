package handlers

import (
	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	svc    service.SVCInterface
}
