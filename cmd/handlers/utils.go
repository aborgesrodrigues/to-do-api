package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aborgesrodrigues/to-do-api/internal/common"
	"github.com/aborgesrodrigues/to-do-api/internal/service"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) *Handler {
	svc, err := service.New(service.Config{Logger: logger})
	if err != nil {
		panic(err)
	}

	return &Handler{
		logger: logger,
		svc:    svc,
	}
}

func writeResponse(w http.ResponseWriter, status int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
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
