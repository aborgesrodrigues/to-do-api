package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const envJWTSecretKey = "JWT_SECRET_KEY"

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

func generateJWT(user *common.User) (string, error) {
	claims := common.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
		CustomClaims: map[string]any{
			"user": user,
		},
	}

	// generate a string using claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecretKey := viper.GetString(envJWTSecretKey)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
