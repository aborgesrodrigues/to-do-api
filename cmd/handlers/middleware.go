package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

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
		handler.logger.Info("IdMiddleware")
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
		logger := handler.logger.With(
			zap.String("from", r.RemoteAddr),
			zap.String("protocol", r.Proto),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		start := time.Now()
		logger.Info("Inbound HTTP request: received",
			zap.Int64("contentLength", r.ContentLength),
		)
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
		logger.Info("Inbound HTTP request: completed",
			zap.Int("statusCode", ww.Status()),
			zap.String("statusText", http.StatusText(ww.Status())),
			zap.Duration("responseTime", time.Since(start).Round(time.Millisecond)),
			zap.Int("bytesWritten", ww.BytesWritten()),
		)

		bReq, err := json.Marshal(getRequestMetadata(clone))
		if err != nil {
			logger.Error("error marshaling request", zap.Error(err))
		}
		logger.Info("request", zap.ByteString("data", bReq))

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

		bRes, err := json.Marshal(getResponseMetadata(respCopy))
		if err != nil {
			logger.Error("error marshaling request", zap.Error(err))
		}
		logger.Info("response", zap.ByteString("data", bRes))
	})
}
