package handlers

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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
		handler.logger.Debug(taskId)

		if userId == "" && taskId == "" {
			handler.logger.Error("User id not passed.")
			writeResponse(rw, http.StatusBadRequest, "User id not passed.")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

// LoggerMiddleware will log request and response for each call
func (handler *Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// copy request body
		var reqBuf bytes.Buffer
		tee := io.TeeReader(r.Body, &reqBuf)
		r.Body = ioutil.NopCloser(tee)

		//copy response body
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
		var resBuf bytes.Buffer
		ww.Tee(&resBuf)

		next.ServeHTTP(ww, r)

		// log request
		clone := r.Clone(context.Background())
		if reqBuf.Len() > 0 {
			clone.Body = ioutil.NopCloser(&reqBuf)
		} else {
			// Server requests are always non-nil, which causes httputil.DumpRequest to
			// write 'null' as the body representation when in fact there was none.
			clone.Body = nil
		}

		handler.logger.Info("request", zap.Any("data", getRequestMetadata(clone)))

		// log response
		respCopy := &http.Response{
			Body:          ioutil.NopCloser(&resBuf),
			ContentLength: int64(ww.BytesWritten()),
			Header:        ww.Header(),
			Proto:         r.Proto,
			ProtoMajor:    r.ProtoMajor,
			ProtoMinor:    r.ProtoMinor,
			Request:       r,
			StatusCode:    ww.Status(),
		}

		handler.logger.Info("response", zap.Any("data", getResponseMetadata(respCopy)))
	})
}
