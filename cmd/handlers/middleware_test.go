package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIdMiddleware(t *testing.T) {
	tests := map[string]struct {
		path       string
		wantStatus int
		wantID     string
	}{
		"no user_id": {
			path:       "/users/",
			wantStatus: http.StatusNotFound,
			wantID:     "",
		},
		"int user_id": {
			path:       "/users/42",
			wantStatus: http.StatusOK,
			wantID:     "42",
		},
		"string user_id": {
			path:       "/users/acb",
			wantStatus: http.StatusOK,
			wantID:     "acb",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			hdl := &Handler{
				Logger: zap.NewNop(),
			}
			var gotID string

			probe := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotID, _ = r.Context().Value(idCtx).(string)
				w.WriteHeader(http.StatusOK)
			})

			r := chi.NewRouter()
			r.Route("/users", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Use(hdl.IdMiddleware)
					r.Get("/", probe)
				})
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantID, gotID)

		})
	}
}

func TestValidateUserMiddleware(t *testing.T) {
	tests := map[string]struct {
		path       string
		userID     string
		wantStatus int
	}{
		"success": {
			path:       "/users/1",
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		"different user": {
			path:       "/users/99",
			userID:     "1",
			wantStatus: http.StatusForbidden,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			hdl := &Handler{
				Logger: zap.NewNop(),
			}

			probe := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			r := chi.NewRouter()
			r.Route("/users/{id}", func(r chi.Router) {
				// a ordem importa: Validate depende do que o IdMiddleware coloca no context
				r.Use(hdl.IdMiddleware)
				r.Use(hdl.ValidateUserMiddleware)
				r.Get("/", probe)
			})

			ctx := context.WithValue(t.Context(), claimsKey,
				&common.Claims{
					UserID: tt.userID,
				},
			)

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil).WithContext(ctx))

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
