package logging

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/clbanning/mxj/v2"
	"go.uber.org/zap"
)

const redacted = "[Redacted]"

// RedactionOptions holds redaction options for http dumps.
type RedactionOptions struct {
	// Use canonical header names.
	RedactHeaders []string
	// Key or dot path notation.
	RedactBodyKeys []string
}

func redactRequest(logger *zap.Logger, r *http.Request, ro RedactionOptions) (err error) {
	if len(ro.RedactBodyKeys) == 0 && len(ro.RedactHeaders) == 0 {
		return nil
	}
	if r == nil {
		logger.Warn("request was nil")
		return nil
	}

	if r.Body != nil {
		i, redacted, err := redactBody(logger, ro, r.Body)
		if err != nil {
			return err
		}
		r.Body = redacted
		r.ContentLength = i
	}

	r.Header = redactHeaders(ro, &r.Header)

	return nil
}

func redactResponse(logger *zap.Logger, r *http.Response, ro RedactionOptions) error {
	if len(ro.RedactBodyKeys) == 0 && len(ro.RedactHeaders) == 0 {
		return nil
	}
	if r == nil {
		logger.Warn("response was nil")
		return nil
	}

	i, redactedBody, err := redactBody(logger, ro, r.Body)
	if err != nil {
		return err
	}
	if i > 0 {
		r.Body = redactedBody
		r.ContentLength = i
	} else {
		// HTTPClient guarantees that response body is non-nil, which causes httputil.DumpResponse
		// to write 'null' as the response body representation even if none was present.
		r.Body = nil
	}

	r.Header = redactHeaders(ro, &r.Header)

	return nil
}

// this function is used for AWS lambda gateway responses
func redactAPIGatewayResponse(logger *zap.Logger, r *events.APIGatewayProxyResponse, ro RedactionOptions) error {
	if len(ro.RedactBodyKeys) == 0 && len(ro.RedactHeaders) == 0 {
		return nil
	}
	if r == nil {
		logger.Warn("response was nil")
		return nil
	}

	i, redactedBody, err := redactBody(logger, ro, strings.NewReader(r.Body))
	if err != nil {
		return err
	}
	if i > 0 {
		rBytes, err := io.ReadAll(redactedBody)
		if err != nil {
			logger.Error("Error converting redacted body")
		}
		r.Body = string(rBytes)
	} else {
		// HTTPClient guarantees that response body is non-nil, which causes httputil.DumpResponse
		// to write 'null' as the response body representation even if none was present.
		r.Body = ""
	}

	r.Headers = redactAPIGatewayHeaders(ro, r.Headers)

	return nil
}

// returns a redacted copy of headers in an API Gateway proxy response used for AWS lambdas
func redactAPIGatewayHeaders(ro RedactionOptions, headers map[string]string) map[string]string {
	headersClone := make(map[string]string)

	for _, redactedHeader := range ro.RedactHeaders {
		if _, ok := headers[redactedHeader]; ok {
			headersClone[redactedHeader] = redacted
		}
	}

	for headerKey, headerVal := range headers {
		if _, ok := headersClone[headerKey]; !ok {
			headersClone[headerKey] = headerVal
		}
	}

	return headersClone
}

// returns a redacted copy of headers
func redactHeaders(ro RedactionOptions, headers *http.Header) http.Header {
	h := headers.Clone()
	for _, header := range ro.RedactHeaders {
		// Only redact header if present.
		if h.Get(header) != "" {
			h.Set(header, redacted)
		}
	}
	return h
}

// Should be called before redactHeaders() in case someone weirdly wanted to redact "Content-Type"
func redactBody(
	logger *zap.Logger,
	ro RedactionOptions,
	body io.Reader,
) (bytesWritten int64, redacted io.ReadCloser, err error) {
	// TODO: avoid doing anything if body == nil or body redaction config is empty?

	// TODO: Prioritize parsing as JSON vs XML based on MIME type if provided. Otherwise default to
	// trying JSON first, as it's more common for us.

	bodyBytes, err := ioutil.ReadAll(body)
	// io.EOF -> there was no body, which is okay
	if err != nil && err != io.EOF {
		return 0, nil, fmt.Errorf("failed to read body: %w", err)
	}

	// Try JSON
	m, err := mxj.NewMapJson(bodyBytes)
	if err == nil {
		if err := redactBodyMap(logger, ro.RedactBodyKeys, &m); err != nil {
			return 0, nil, fmt.Errorf("failed to redact body: %w", err)
		}
		var buf bytes.Buffer
		if err := m.JsonWriter(&buf, true); err != nil {
			return 0, nil, fmt.Errorf("failed to write json: %w", err)
		}
		return int64(buf.Len()), ioutil.NopCloser(&buf), nil
	}
	logger.Debug("Unable to parse body as JSON.", zap.Error(err))

	// Try XML
	// FIXME: throws an error when document has a out-of-root process instruction, directive, or comment.
	ms, err := mxj.NewMapXmlSeq(bodyBytes)
	if err == nil {
		mMap := mxj.Map(ms)
		if err := redactBodyMap(logger, ro.RedactBodyKeys, &mMap); err != nil {
			return 0, nil, fmt.Errorf("failed to redact body: %w", err)
		}
		var buf bytes.Buffer
		if err := mMap.XmlWriter(&buf); err != nil {
			return 0, nil, fmt.Errorf("failed to write xml: %w", err)
		}
		return int64(buf.Len()), ioutil.NopCloser(&buf), nil
	}
	logger.Debug("Unable to parse body as XML.", zap.Error(err))

	// We only support redaction of JSON and XML, currently. So otherwise, just
	// return a copy of the original stream.
	buf := bytes.NewBuffer(bodyBytes)
	return int64(buf.Len()), ioutil.NopCloser(buf), nil
}

func redactBodyMap(logger *zap.Logger, keysAndPaths []string, m *mxj.Map) error {
	updateValuesForPath := func(key string, path string) error {
		// TODO: this may assume string
		i, err := m.UpdateValuesForPath(key+":"+redacted, path)
		if err != nil {
			logger.Error("Failed to redact body map node.",
				zap.Error(err),
				zap.String("key", key),
				zap.String("path", path),
			)
			// Not worried about exiting early here since we are passing paths generated by
			// PathsForKey(); really shouldn't get an error.
			return err
		}
		logger.Debug("Redacted body map node.",
			zap.Int("occurances", i),
			zap.String("key", key),
			zap.String("path", path),
			zap.String("sigil", redacted),
		)
		return nil
	}
	for _, keyOrPath := range keysAndPaths {
		if strings.Contains(keyOrPath, ".") {
			// It's a path.
			split := strings.Split(keyOrPath, ".")
			key := split[len(split)-1]
			if err := updateValuesForPath(key, keyOrPath); err != nil {
				return err
			}
		}
		// It's a key.
		for _, path := range m.PathsForKey(keyOrPath) {
			if err := updateValuesForPath(keyOrPath, path); err != nil {
				return err
			}
		}
	}
	return nil
}
