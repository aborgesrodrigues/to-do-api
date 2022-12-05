package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ctxKey string

const (
	userIdCtx = ctxKey("userId")
	taskIdCtx = ctxKey("taskId")
)

func (handler *Handler) IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, string(userIdCtx))
		if userId != "" {
			r = r.WithContext(context.WithValue(r.Context(), userIdCtx, userId))
		}
		taskId := chi.URLParam(r, string(taskIdCtx))
		if taskId != "" {
			r = r.WithContext(context.WithValue(r.Context(), taskIdCtx, taskId))
		}
		handler.Logger.Debug(taskId)

		if userId == "" && taskId == "" {
			handler.Logger.Error("User id not passed.")
			writeResponse(rw, http.StatusBadRequest, "User id not passed.")
			return
		}

		next.ServeHTTP(rw, r)
	})
}
