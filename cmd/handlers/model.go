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

// AccessLoggerOptions holds options for constructing the AccessLogger middleware.
type AccessLoggerOptions struct {
	HTTPAuditLogger *logging.HTTPAuditLogger
	// Logger to use in any case where a logger cannot be resolved automatically from
	// request context.
	DefaultLogger *zap.Logger
}
