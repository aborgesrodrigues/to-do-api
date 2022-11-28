package service

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ctxKey string

const (
	userIdCtx = ctxKey("userId")
	taskIdCtx = ctxKey("taskId")
)

func (s *Service) IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, string(userIdCtx))
		if userId != "" {
			r = r.WithContext(context.WithValue(r.Context(), userIdCtx, userId))
		}
		taskId := chi.URLParam(r, string(taskIdCtx))
		if taskId != "" {
			r = r.WithContext(context.WithValue(r.Context(), taskIdCtx, taskId))
		}
		s.Logger.Debug(taskId)

		if userId == "" && taskId == "" {
			s.Logger.Error("User id not passed.")
			writeResponse(rw, http.StatusBadRequest, "User id not passed.")
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func (s *Service) LoggerMiddleware(next http.Handler) http.Handler {
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

		s.Logger.Info("request", zap.Any("data", getRequestMetadata(clone)))

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

		s.Logger.Info("response", zap.Any("data", getResponseMetadata(respCopy)))
	})
}

func getRequestMetadata(req *http.Request) []common.Metadata {
	var m []common.Metadata
	m = append(m,
		common.Metadata{Name: "host", Value: req.Host}, // TODO: test client vs server request behavior
		common.Metadata{Name: "hostname", Value: req.URL.Hostname()},
		common.Metadata{Name: "method", Value: req.Method},
		common.Metadata{Name: "path", Value: req.URL.Path},
		common.Metadata{Name: "protocol", Value: req.Proto},
		common.Metadata{Name: "query", Value: req.URL.Query().Encode()},
		common.Metadata{Name: "fragment", Value: req.URL.Fragment},
		common.Metadata{Name: "headers", Value: req.Header},
	)
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			m = append(m, common.Metadata{Name: "bodyReadError", Value: err.Error()})
		}
		m = append(m, common.Metadata{Name: "body", Value: string(body)})
	}
	return m
}

func getResponseMetadata(res *http.Response) []common.Metadata {
	var m []common.Metadata
	m = append(m,
		common.Metadata{Name: "protocol", Value: res.Proto},
		common.Metadata{Name: "requestHost", Value: res.Request.Host},
		common.Metadata{Name: "requestHostname", Value: res.Request.URL.Hostname()},
		common.Metadata{Name: "requestMethod", Value: res.Request.Method},
		common.Metadata{Name: "requestPath", Value: res.Request.URL.Path},
		common.Metadata{Name: "requestProtocol", Value: res.Request.Proto},
		common.Metadata{Name: "status", Value: res.Status},
		common.Metadata{Name: "statusCode", Value: res.StatusCode},
		common.Metadata{Name: "headers", Value: res.Header},
	)
	if res.Body != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			m = append(m, common.Metadata{Name: "bodyReadError", Value: err.Error()})
		}
		m = append(m, common.Metadata{Name: "body", Value: string(body)})
	}
	return m
}
