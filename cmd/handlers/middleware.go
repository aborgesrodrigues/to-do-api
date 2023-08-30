package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/aborgesrodrigues/to-do-api/internal/logging"
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
		handler.Logger.Debug(taskId)

		if userId == "" && taskId == "" {
			handler.Logger.Error("User id not passed.")
			writeResponse(rw, http.StatusBadRequest, "User id not passed.")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

// AccessLogger returns a Chi middleware which writes access logs using zap.
// This requires RequestCorrelationLogger having been called earlier in the middleware chain.
func AccessLogger(opt AccessLoggerOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// Need to make a copy of ctx before next.ServeHTTP is called.
			ctx := req.Context()
			logger, err := logging.LoggerFromContext(req.Context())
			// TODO: lazy-instantiated default logger instead (internal to lib)
			if err != nil {
				if opt.DefaultLogger != nil {
					logger = opt.DefaultLogger
				} else {
					logger = zap.NewNop()
				}
			}
			logger = logger.With(
				zap.String("from", req.RemoteAddr),
				zap.String("protocol", req.Proto),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
			)
			start := time.Now()
			logger.Info("Inbound HTTP request: received",
				zap.Int64("contentLength", req.ContentLength),
			)
			ww := middleware.NewWrapResponseWriter(res, req.ProtoMajor)
			var reqBuf bytes.Buffer
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.Logger != nil &&
				!opt.HTTPAuditLogger.Opt.DisableRequestAuditLogs {

				// If audit logging request, we'll need a copy of the request body.
				tee := io.TeeReader(req.Body, &reqBuf)
				req.Body = io.NopCloser(tee)
			}
			var resBuf bytes.Buffer
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.Logger != nil &&
				!opt.HTTPAuditLogger.Opt.DisableResponseAuditLogs {

				// If audit logging response, we'll need a copy of the response body.
				ww.Tee(&resBuf)
			}
			next.ServeHTTP(ww, req)
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.Logger != nil &&
				!opt.HTTPAuditLogger.Opt.DisableRequestAuditLogs {

				// Doing this after ServeHTTP so request body will be written to TeeReader.
				clone := req.Clone(context.Background())
				if reqBuf.Len() > 0 {
					clone.Body = io.NopCloser(&reqBuf)
				} else {
					// Server requests are always non-nil, which causes httputil.DumpRequest to
					// write 'null' as the body representation when in fact there was none.
					clone.Body = nil
				}
				go opt.HTTPAuditLogger.LogUpstreamRequest(ctx, logger, clone)
			}
			logger.Info("Inbound HTTP request: completed",
				zap.Int("statusCode", ww.Status()),
				zap.String("statusText", http.StatusText(ww.Status())),
				zap.Duration("responseTime", time.Since(start).Round(time.Millisecond)),
				zap.Int("bytesWritten", ww.BytesWritten()),
			)
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.Logger != nil &&
				!opt.HTTPAuditLogger.Opt.DisableResponseAuditLogs {

				respCopy := &http.Response{
					Body:          io.NopCloser(&resBuf),
					ContentLength: int64(ww.BytesWritten()),
					Header:        ww.Header(),
					Proto:         req.Proto,
					ProtoMajor:    req.ProtoMajor,
					ProtoMinor:    req.ProtoMinor,
					// TODO: should we clone the request?
					Request:    req,
					StatusCode: ww.Status(),
				}
				go opt.HTTPAuditLogger.LogUpstreamResponse(ctx, logger, respCopy)
			}
		})
	}
}
