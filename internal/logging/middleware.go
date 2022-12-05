package logging

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type contextKey string

const ctxKeyLogger = contextKey("request-logger")

func RequestCorrelationLogger(baseLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			logger := baseLogger

			ctx := ContextWithLogger(req.Context(), logger)

			reqWithCtx := req.WithContext(ctx)
			next.ServeHTTP(res, reqWithCtx)
		})
	}
}

func ContextWithLogger(parent context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(parent, ctxKeyLogger, logger)
}

func LoggerFromContext(ctx context.Context) (*zap.Logger, error) {
	maybeLogger := ctx.Value(ctxKeyLogger)
	if maybeLogger == nil {
		return nil, errors.New("logger not on context")
	}
	logger, ok := maybeLogger.(*zap.Logger)
	if !ok {
		return nil, errors.New("context value not a *zap.Logger")
	}
	return logger, nil
}
