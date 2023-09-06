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

func generateJWT(user *common.User) (string, string, error) {
	jwtSecretKey := viper.GetString(envJWTSecretKey)

	// access token
	accessClaims := common.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Second)),
		},
		Type:   common.AccessTokenType,
		UserID: user.Id,
	}

	// generate a string using claims
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessTokenString, err := accessToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", err
	}

	// refresh token
	refreshClaims := common.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		Type:   common.RefreshTokenType,
		UserID: user.Id,
	}

	// generate a string using claims
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
