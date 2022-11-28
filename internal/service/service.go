package service

import (
	"encoding/json"
	"net/http"

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
		Logger: logger,
		DB:     db,
	}, nil
}

func writeResponse(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

func (svc *Service) HelloWorld(w http.ResponseWriter, r *http.Request) {
	svc.Logger.Info("HelloWorld")

	writeResponse(w, http.StatusOK, map[string]string{
		"message": "Hello World",
	})
}
