package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ctxKey string

const (
	idCtx = ctxKey("Id")
)

func (handler *Handler) IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, string(idCtx))
		if id != "" {
			r = r.WithContext(context.WithValue(r.Context(), idCtx, id))
		}

		handler.Logger.Debug(id)

		if id == "" {
			handler.Logger.Error("Id not passed.")
			writeResponse(rw, http.StatusBadRequest, "Id not passed.")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

var reqLatencyAuditHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "request_latency_audit",
	Help: "Latency of Audit",
	Buckets: []float64{
		.1,
		1,
		25,
		50,
		100,
		150,
		200,
		250,
		500,
		1000,
		2500,
		5000,
		10000,
	},
}, []string{"disable_request_audit_logs", "disable_response_audit_logs", "buffer_size"})

// AccessLogger returns a Chi middleware which writes access logs using zap.
// This requires RequestCorrelationLogger having been called earlier in the middleware chain.
func AccessLogger(opt AccessLoggerOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			startMetric := time.Now()
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

			delay := time.Since(startMetric).Milliseconds()
			labels := prometheus.Labels{
				"disable_request_audit_logs":  strconv.FormatBool(opt.HTTPAuditLogger.Opt.DisableResponseAuditLogs),
				"disable_response_audit_logs": strconv.FormatBool(opt.HTTPAuditLogger.Opt.DisableResponseAuditLogs),
				"buffer_size":                 strconv.Itoa(opt.HTTPAuditLogger.Opt.Config.BufferSize),
			}

			reqLatencyAuditHistogram.With(labels).Observe(float64(delay))
		})
	}
}

func (handler *Handler) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		claims, err := handler.validateJWT(r)
		if err != nil {
			handler.Logger.Error("Error validating JWT")
			writeResponse(rw, http.StatusUnauthorized, err.Error())
			return
		}

		// force use of access token only
		if claims.Type != common.AccessTokenType || claims.UserID == "" {
			handler.Logger.Error("Use of the invalid type of jwt token")
			writeResponse(rw, http.StatusUnauthorized, "You're Unauthorized due to Invalid type of jwt token")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func (handler *Handler) VerifyRefreshJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		claims, err := handler.validateJWT(r)
		if err != nil {
			handler.Logger.Error("Error validating JWT")
			writeResponse(rw, http.StatusUnauthorized, err.Error())
			return
		}

		// force use of access token only
		if claims.Type != common.RefreshTokenType || claims.UserID == "" {
			handler.Logger.Error("Use of the invalid type of jwt token")
			writeResponse(rw, http.StatusUnauthorized, "You're Unauthorized due to Invalid type of jwt token")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func (handler *Handler) validateJWT(r *http.Request) (*common.Claims, error) {
	var id string
	if r.Context().Value(id) != nil {
		id = r.Context().Value(id).(string)
	}

	if headerToken := r.Header.Get("Authorization"); headerToken != "" {
		headerToken = strings.Replace(headerToken, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(headerToken, &common.Claims{}, func(token *jwt.Token) (interface{}, error) {
			jwtSecretKey := viper.GetString(envJWTSecretKey)
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				handler.Logger.Error("Malformed token")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				handler.Logger.Error("Invalid Signature")
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				handler.Logger.Error("Invalid time")
			default:
				handler.Logger.Error("Error parsing token")
			}

			return nil, err
		}

		if !token.Valid {
			handler.Logger.Error("Invalid token")
			return nil, errors.New("invalid token")
		}
		claims, ok := token.Claims.(*common.Claims)
		if !ok {
			handler.Logger.Error("Invalid type of claims")
			return nil, errors.New("invalid type of claims")
		}

		// TODO improve validation of a token for a specific id
		if id != "" && claims.UserID != id {
			handler.Logger.Error("Invalid token for this id")
			return nil, errors.New("invalid token for this id")
		}

		handler.Logger.Info("Valid user", zap.String("user", claims.UserID))

		return claims, nil
	}

	return nil, errors.New("no authorization header informed")
}

var reqLatencyHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "request_latency",
	Help: "Latency of HTTP requests",
	Buckets: []float64{
		.1,
		1,
		25,
		50,
		100,
		150,
		200,
		250,
		500,
		1000,
		2500,
		5000,
		10000,
	},
}, []string{"path", "status"})

var reqCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "request_count",
	Help: "Latency of HTTP requests",
}, []string{"path", "status"})

func (handler *Handler) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		delay := time.Since(start).Milliseconds()
		labels := prometheus.Labels{
			"path":   r.URL.Path,
			"status": strconv.Itoa(ww.Status()),
		}

		reqLatencyHistogram.With(labels).Observe(float64(delay))
		reqCount.With(labels).Inc()
	})
}
