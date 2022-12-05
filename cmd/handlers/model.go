package handlers

import (
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	Logger      *zap.Logger
	AuditLogger *logging.HTTPAuditLogger
	svc         service.SVCInterface
}
