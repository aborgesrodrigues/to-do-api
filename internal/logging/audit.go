package logging

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aborgesrodrigues/to-do-api/internal/audit"
)

// HTTPAuditLogger facilitates logging HTTP requests and responses
type HTTPAuditLogger struct {
	Opt      HTTPAuditLogOptions
	Logger   *audit.Logger
	wg       *sync.WaitGroup
	once     sync.Once
	isClosed bool
}

// Close closes the underlying audit.Logger and confirms there are no outstanding tasks in wait group.
// Client should call Close before its process exits.
func (l *HTTPAuditLogger) Close() {
	// wait for up to 2 secs
	l.once.Do(func() {
		l.isClosed = true
		waitTimeout(l.wg, 2*time.Second)
		l.Logger.Close()
	})
}

// HTTPAuditLogOptions provides options to NewHTTPAuditLogger()
type HTTPAuditLogOptions struct {
	AuditPathParams          []string
	DisableRequestAuditLogs  bool
	DisableResponseAuditLogs bool

	audit.Config
	Writer audit.EventWriter
}

// NewHTTPAuditLogger creates a new HTTPAuditLogger.
func NewHTTPAuditLogger(opt HTTPAuditLogOptions) (*HTTPAuditLogger, error) {
	al, err := audit.NewLogger(opt.Config, opt.Writer)
	if err != nil {
		return nil, err
	}

	al.SetOnError(func(ctx context.Context, err error) {
		logger, logerr := LoggerFromContext(ctx)
		if logerr == nil {
			logger.Info("Failed to write audit log.", zap.Error(err))
		}
	})
	al.SetOnWrite(func(ctx context.Context, out map[string]interface{}) {
		logger, err := LoggerFromContext(ctx)
		if err == nil {
			fields := make([]zapcore.Field, 0, len(out))
			for i := range out {
				fields = append(fields, zap.Any(i, out[i]))
			}
			logger.Info("Wrote audit log.", fields...)
		}
	})

	return &HTTPAuditLogger{
		Logger: al,
		Opt:    opt,
		wg:     &sync.WaitGroup{},
	}, nil
}

// e.g. /foo/123/bar/456/baz -> /foo/~/bar/~/baz
func (l *HTTPAuditLogger) normalizePath(path string) string {
	split := func(p string) []string {
		return strings.Split(strings.Trim(p, "/"), "/")
	}
PathLoop:
	for _, p := range l.Opt.AuditPathParams {
		known := split(p)
		unknown := split(path)
		if len(known) != len(unknown) {
			continue
		}
		for i := range known {
			// TODO: trim whitespace
			if strings.HasPrefix(known[i], "{") && strings.HasSuffix(known[i], "}") {
				unknown[i] = "~"
			} else if known[i] != unknown[i] {
				continue PathLoop
			}
		}
		return "/" + strings.Join(unknown, "/")
	}
	return path
}

func (l *HTTPAuditLogger) makeAuditID(req *http.Request, downstream bool, response bool) string {
	builder := strings.Builder{}
	if downstream {
		builder.WriteString("out/")
		builder.WriteString(req.URL.Hostname())
	} else {
		builder.WriteString("in")
	}
	path := l.normalizePath(req.URL.EscapedPath())
	builder.WriteString(path)
	if !strings.HasSuffix(path, "/") {
		builder.WriteString("/")
	}
	builder.WriteString(req.Method)
	if response {
		builder.WriteString("/response")
	} else {
		builder.WriteString("/request")
	}
	return builder.String()
}

func (l *HTTPAuditLogger) makeRequestID(req *http.Request, downstream bool) string {
	return l.makeAuditID(req, downstream, false)
}

func (l *HTTPAuditLogger) makeResponseID(resp *http.Response, downstream bool) string {
	return l.makeAuditID(resp.Request, downstream, true)
}

func getRequestMetadata(ctx context.Context, req *http.Request) []audit.Metadata {
	var m []audit.Metadata
	m = append(m,
		audit.Metadata{Name: "host", Value: req.Host}, // TODO: test client vs server request behavior
		audit.Metadata{Name: "hostname", Value: req.URL.Hostname()},
		audit.Metadata{Name: "method", Value: req.Method},
		audit.Metadata{Name: "path", Value: req.URL.Path},
		audit.Metadata{Name: "protocol", Value: req.Proto},
		audit.Metadata{Name: "query", Value: req.URL.Query().Encode()},
		audit.Metadata{Name: "fragment", Value: req.URL.Fragment},
		audit.Metadata{Name: "headers", Value: req.Header},
	)
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			m = append(m, audit.Metadata{Name: "bodyReadError", Value: err.Error()})
		}
		m = append(m, audit.Metadata{Name: "body", Value: string(body)})
	}

	return m
}

func getResponseMetadata(ctx context.Context, res *http.Response) []audit.Metadata {
	var m []audit.Metadata
	m = append(m,
		audit.Metadata{Name: "protocol", Value: res.Proto},
		audit.Metadata{Name: "requestHost", Value: res.Request.Host},
		audit.Metadata{Name: "requestHostname", Value: res.Request.URL.Hostname()},
		audit.Metadata{Name: "requestMethod", Value: res.Request.Method},
		audit.Metadata{Name: "requestPath", Value: res.Request.URL.Path},
		audit.Metadata{Name: "requestProtocol", Value: res.Request.Proto},
		audit.Metadata{Name: "status", Value: res.Status},
		audit.Metadata{Name: "statusCode", Value: res.StatusCode},
		audit.Metadata{Name: "headers", Value: res.Header},
	)
	if res.Body != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			m = append(m, audit.Metadata{Name: "bodyReadError", Value: err.Error()})
		}
		m = append(m, audit.Metadata{Name: "body", Value: string(body)})
	}

	return m
}

// call method wrapped with executeAsync if you want async execution
// write audit log for a request from upstream
// req should be safe for mutation
func (l *HTTPAuditLogger) LogUpstreamRequest(
	ctx context.Context,
	logger *zap.Logger,
	req *http.Request,
) {
	if l.Logger == nil {
		logger.Error("Unable to write upstream request audit log. Audit logger not provided.")
		return
	}

	id := l.makeRequestID(req, false)
	metadata := getRequestMetadata(ctx, req)
	l.Logger.Write(ctx, id, metadata...)
}

// call method wrapped with executeAsync if you want async execution
// write audit log for a response to upstream
// resp should be safe for mutation
func (l *HTTPAuditLogger) LogUpstreamResponse(
	ctx context.Context,
	logger *zap.Logger,
	resp *http.Response,
) {
	if l.Logger == nil {
		logger.Error("Unable to write upstream response audit log. Audit logger not provided.")
		return
	}

	id := l.makeResponseID(resp, false)
	metadata := getResponseMetadata(ctx, resp)
	l.Logger.Write(ctx, id, metadata...)
}

// waitTimeout waits for either a duration to elapse or a wait group to be done before returning, whichever happens first
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
