package logging

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// AccessLoggerOptions holds options for constructing the AccessLogger middleware.
type AccessLoggerOptions struct {
	HTTPAuditLogger *HTTPAuditLogger
	// Logger to use in any case where a logger cannot be resolved automatically from
	// request context.
	DefaultLogger *zap.Logger
}

// AccessLogger returns a Chi middleware which writes access logs using zap.
// This requires RequestCorrelationLogger having been called earlier in the middleware chain.
func AccessLogger(opt AccessLoggerOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// Need to make a copy of ctx before next.ServeHTTP is called.
			ctx := req.Context()
			logger, err := LoggerFromContext(req.Context())
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
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.logger != nil &&
				!opt.HTTPAuditLogger.opt.DisableRequestAuditLogs {

				// If audit logging request, we'll need a copy of the request body.
				tee := io.TeeReader(req.Body, &reqBuf)
				req.Body = ioutil.NopCloser(tee)
			}
			var resBuf bytes.Buffer
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.logger != nil &&
				!opt.HTTPAuditLogger.opt.DisableResponseAuditLogs {

				// If audit logging response, we'll need a copy of the response body.
				ww.Tee(&resBuf)
			}
			next.ServeHTTP(ww, req)
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.logger != nil &&
				!opt.HTTPAuditLogger.opt.DisableRequestAuditLogs {

				// Doing this after ServeHTTP so request body will be written to TeeReader.
				clone := req.Clone(context.Background())
				if reqBuf.Len() > 0 {
					clone.Body = ioutil.NopCloser(&reqBuf)
				} else {
					// Server requests are always non-nil, which causes httputil.DumpRequest to
					// write 'null' as the body representation when in fact there was none.
					clone.Body = nil
				}
				go opt.HTTPAuditLogger.logUpstreamRequest(ctx, logger, clone)
			}
			logger.Info("Inbound HTTP request: completed",
				zap.Int("statusCode", ww.Status()),
				zap.String("statusText", http.StatusText(ww.Status())),
				zap.Duration("responseTime", time.Since(start).Round(time.Millisecond)),
				zap.Int("bytesWritten", ww.BytesWritten()),
			)
			if opt.HTTPAuditLogger != nil && opt.HTTPAuditLogger.logger != nil &&
				!opt.HTTPAuditLogger.opt.DisableResponseAuditLogs {

				respCopy := &http.Response{
					Body:          ioutil.NopCloser(&resBuf),
					ContentLength: int64(ww.BytesWritten()),
					Header:        ww.Header(),
					Proto:         req.Proto,
					ProtoMajor:    req.ProtoMajor,
					ProtoMinor:    req.ProtoMinor,
					// TODO: should we clone the request?
					Request:    req,
					StatusCode: ww.Status(),
				}
				go opt.HTTPAuditLogger.logUpstreamResponse(ctx, logger, respCopy)
			}
		})
	}
}
