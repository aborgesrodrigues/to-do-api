package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"go.uber.org/zap"
)

func New(logger *zap.Logger, auditLogger *logging.HTTPAuditLogger) *Handler {
	svc, err := service.New(service.Config{Logger: logger})
	if err != nil {
		panic(err)
	}
	logger.Info("handler created")
	return &Handler{
		Logger:      logger,
		AuditLogger: auditLogger,
		svc:         svc,
	}
}

func writeResponse(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}
