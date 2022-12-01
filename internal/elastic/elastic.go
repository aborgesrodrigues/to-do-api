package elastic

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds the values needed to initialize the elasticSink
type Config struct {
	Level         *zap.AtomicLevel
	Client        *Client
	Endpoint      string
	Username      *string
	Password      *string
	Index         string
	CaCert        *string
	FlushInterval *time.Duration // The periodic flush interval
	FlushBytes    *int           // The flush threshold in bytes
	NumWorkers    *int           // The number of worker goroutines
}

type elasticSink struct {
	cfg         Config
	Client      *Client
	Index       string
	BulkIndexer *esutil.BulkIndexer
}

// Client wraps the elasticsearch client
type Client struct {
	*elasticsearch.Client
}

// Sink returns the zap.Sink for elasticsearch
func Sink(u *url.URL) (zap.Sink, error) {
	err := validateURL(u)
	if err != nil {
		return nil, err
	}
	pwd, _ := u.User.Password()
	caFile := u.Query().Get("caCert")
	user := u.User.Username()

	return getElasticSink(Config{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username: &user,
		Password: &pwd,
		Index:    u.Path,
		CaCert:   &caFile,
	})
}

func getScheme(u *url.URL) (*string, error) {
	scheme := u.Scheme

	if scheme != "https" && scheme != "http" {
		return nil, fmt.Errorf("elastic-logger: scheme must be http or https. got: %s", scheme)
	}

	return &scheme, nil
}

func validateCA(cas string, tp *http.Transport) (*http.Transport, error) {
	var ca []byte

	if cas != "" {
		var err error
		caString, err := ioutil.ReadFile(cas)
		if err != nil {
			return nil, err
		}

		ca = []byte(caString)

		if tp.TLSClientConfig.RootCAs, err = x509.SystemCertPool(); err != nil {
			return nil, fmt.Errorf("elastic-logger: failed to add CA to pool: %w", err)
		}

		if ok := tp.TLSClientConfig.RootCAs.AppendCertsFromPEM(ca); !ok {
			return nil, fmt.Errorf("elastic-logger: failed to add PEM from CA: %w", err)
		}
	}

	return tp, nil
}

func validateURL(u *url.URL) error {
	if u.Path == "" {
		return errors.New("elastic-logger: path cannot be empty; should contain the index in elasticsearch. e.g. https://user:pwd@localhost:9200/your_index")
	}

	match, err := regexp.MatchString(`^\/(.*)\/{1,}`, u.Path)
	if match {
		return errors.New("elastic-logger: index cannot contain more than a single path")
	}
	if err != nil {
		return fmt.Errorf("elastic-logger: failed to test index against path regex. err: %w", err)
	}

	_, ok := u.User.Password()
	if !ok {
		return errors.New("elastic-logger: password cannot be nil")
	}

	return nil
}

func getElasticSink(cfg Config) (*elasticSink, error) {
	var (
		client *Client
		err    error
	)

	// if client didn't bring their own Elasticsearch client, create our own
	if cfg.Client == nil {
		client, _, err = NewElasticClient(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise use the one provided
		client = cfg.Client
	}

	bulkIndexer, err := NewBulkIndexer(cfg, client)
	if err != nil {
		return nil, err
	}

	return &elasticSink{
		cfg:         cfg,
		Client:      client,
		BulkIndexer: bulkIndexer,
		Index:       cfg.Index,
	}, nil
}

// getElasticConfig takes the elastic-logger.Config and converts it into an elasticsearch.Config
func getElasticConfig(cfg Config) (*elasticsearch.Config, error) {
	var tp *http.Transport
	var err error

	tp = http.DefaultTransport.(*http.Transport).Clone()

	if cfg.CaCert != nil {
		tp, err = validateCA(*cfg.CaCert, tp)
		if err != nil {
			return nil, err
		}
	}

	clientConfig := elasticsearch.Config{
		Addresses: []string{cfg.Endpoint},
		Transport: tp,
	}

	if cfg.Username != nil {
		clientConfig.Username = *cfg.Username
	}

	if cfg.Password != nil {
		clientConfig.Password = *cfg.Password
	}

	return &clientConfig, nil
}

// NewElasticClient instantiates a new Elasticsearch client
func NewElasticClient(cfg Config) (*Client, *esapi.Response, error) {
	config, err := getElasticConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	client, err := elasticsearch.NewClient(*config)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove. Since v7.14.0, the elasticsearch client itself makes a check
	// which is redundant with this one.
	status, err := client.Info()
	if err != nil {
		return nil, nil, fmt.Errorf("elastic-logger: unable to get elastic info: %w", err)
	}

	exists, err := client.Indices.Exists([]string{strings.TrimLeft(cfg.Index, "/")})
	if err != nil {
		return nil, nil, err
	}
	if exists.StatusCode != 200 {
		var (
			buf    bytes.Buffer
			reader *bytes.Reader
		)

		b := []byte(`{ "mappings": { "properties": { "@timestamp": { "type": "date" } } } }`)
		buf.Grow(len(b))
		buf.Write(b)

		reader = bytes.NewReader(buf.Bytes())
		req := esapi.IndicesCreateRequest{
			Index: strings.TrimLeft(cfg.Index, "/"),
			Body:  reader,
		}
		createIndex, err := req.Do(context.Background(), client)
		if err != nil {
			return nil, nil, err
		}

		if createIndex.IsError() {
			return nil, nil, fmt.Errorf(
				"elastic-logger: unable to create elasticsearch index '%s': %s",
				cfg.Index,
				// keeping newlines from the elasticsearch response out of the error message
				strings.TrimSpace(createIndex.String()),
			)
		}
	}

	elasticClient := Client{
		client,
	}

	return &elasticClient, status, nil
}

// elasticErrorHandler parses an Elasticsearch response when errors occur and returns a Golang error
// with that information
func elasticErrorHandler(res esutil.BulkIndexerResponseItem) error {
	js, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("elastic-logger: failed to marshal response body: %w", err)
	}
	errStr := fmt.Sprintf("elasticsearch error: [%d] %s: %s; response body: %s",
		res.Status,
		res.Error.Type,
		res.Error.Reason,
		string(js),
	)
	return errors.New(errStr)
}

func (e *elasticSink) Write(b []byte) (n int, err error) {
	var (
		buf    bytes.Buffer
		reader *bytes.Reader
	)

	// add new line here to prevent sending bunch of logs as one HUGE log
	b = append(b, "\n"...)
	buf.Grow(len(b))
	buf.Write(b)

	reader = bytes.NewReader(buf.Bytes())

	err = (*e.BulkIndexer).Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action: "index", // Action field configures the operation to perform (index, create, delete, update)
			Body:   reader,
			OnFailure: func(ctx context.Context, req esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				fmt.Println(elasticErrorHandler(res))
			},
		},
	)
	if err != nil {
		return len(b), err
	}

	return len(b), nil
}

// Close is required to implement zap.Sink func Close
// https://github.com/uber-go/zap/blob/master/sink.go#L63
func (e *elasticSink) Close() error {
	return (*e.BulkIndexer).Close(context.Background())
}

// Sync is required to implement zap.Sink func Sync
func (e *elasticSink) Sync() error {
	closeErr := (*e.BulkIndexer).Close(context.Background())
	if closeErr != nil {
		return closeErr
	}

	bulkIndexer, bulkErr := NewBulkIndexer(e.cfg, e.Client)
	if bulkErr != nil {
		return bulkErr
	}
	e.BulkIndexer = bulkIndexer

	return nil
}

/*NewBulkIndexer instantiates a new Bulk client
Default bulk config values:
	1) NumWorkers = 1
	2) FlushBytes = 5e+6
	3) FlushInterval = 5 * time.Second*/
func NewBulkIndexer(cfg Config, elasticClient *Client) (*esutil.BulkIndexer, error) {
	var (
		numWorkers        = runtime.NumCPU()
		flushBytes    int = 5e+6            // Default 5mb
		flushInterval     = 5 * time.Second // Default 5 seconds
	)

	// overwrite default configs
	if cfg.NumWorkers != nil {
		numWorkers = *cfg.NumWorkers
	}
	if cfg.FlushBytes != nil {
		flushBytes = *cfg.FlushBytes
	}
	if cfg.FlushInterval != nil {
		flushInterval = *cfg.FlushInterval
	}

	bulkCfg := esutil.BulkIndexerConfig{
		Index:         strings.TrimLeft(cfg.Index, "/"), // Default index name
		Client:        elasticClient.Client,
		NumWorkers:    numWorkers,
		FlushBytes:    flushBytes,
		FlushInterval: flushInterval,
	}

	bi, err := esutil.NewBulkIndexer(bulkCfg)
	if err != nil {
		return nil, err
	}

	return &bi, nil
}

// NewCore registers the Elasticsearch sink and returns an ECS compatible zapcore.Core
func (cfg *Config) NewCore() (zapcore.Core, error) {
	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	scheme, err := getScheme(u)
	if err != nil {
		return nil, err
	}
	zap.RegisterSink(*scheme, Sink)

	if cfg.Username != nil && cfg.Password != nil {
		u.User = url.UserPassword(*cfg.Username, *cfg.Password)
		u.Path = cfg.Index
	}

	if cfg.CaCert != nil {
		q := u.Query()
		q.Set("caCert", *cfg.CaCert)
		u.RawQuery = q.Encode()
	}

	sink, err := getElasticSink(*cfg)
	if err != nil {
		return nil, err
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	return ecszap.NewCore(encoderConfig, sink, cfg.Level), nil
}
