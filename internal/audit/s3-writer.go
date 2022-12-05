package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Writer struct {
	config   S3Config
	s3Client *s3.S3
}

// S3Config contains configuration details for creating an S3 event writer.
type S3Config struct {
	// An optional endpoint URL (hostname only or fully qualified URI)
	// that overrides the default generated endpoint for a client. Set this
	// to `nil` to use the default generated endpoint.
	//
	// Note: You must still provide a `Region` value when specifying an
	// endpoint for a client.
	Endpoint  *string
	Region    string
	Bucket    string
	Directory string
}

// NewS3Writer creates a new audit event writer that outputs to an S3 bucket.
// All parameters must be non-empty strings, except for the last parameter (Directory) which is optional.
func NewS3Writer(config S3Config) (EventWriter, error) {
	missing := func(k, v string) error {
		if v == "" {
			return fmt.Errorf("missing required config value: %s", k)
		}
		return nil
	}
	if err := missing("Region", config.Region); err != nil {
		return nil, err
	}
	if err := missing("Bucket", config.Bucket); err != nil {
		return nil, err
	}
	if len(config.Directory) > 0 && !strings.HasSuffix(config.Directory, "/") {
		config.Directory = config.Directory + "/"
	}

	// Force S3 path style when AWS endpoint is overridden as providers like
	// localstack (e.g. for local development) may not support the default
	// host prefix pattern.
	s3ForcePathStyle := config.Endpoint != nil

	awsConfig := aws.Config{
		S3ForcePathStyle: aws.Bool(s3ForcePathStyle),
		Endpoint:         config.Endpoint,
		Region:           aws.String(config.Region),
	}

	// Initiate new aws session. Based on aws docs, a session should be cached and
	// reused https://docs.aws.amazon.com/sdk-for-go/api/aws/session
	session, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, err
	}

	// Create new S3 Client and store it on s3Writer
	s3client := s3.New(session)

	return &s3Writer{
		config:   config,
		s3Client: s3client,
	}, nil
}

func (w *s3Writer) uploadToS3(ctx context.Context, identifier string, buf *bytes.Buffer, timestamp time.Time) (map[string]interface{}, error) {
	key := w.config.Directory + getDateDirectory(timestamp) + identifier
	put, err := w.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(w.config.Bucket),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(buf.Bytes()),
		ContentLength:        aws.Int64(int64(buf.Len())),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{}, 3)
	out["region"] = w.config.Region
	out["location"] = fmt.Sprintf("s3://%s/%s", w.config.Bucket, key)
	if put.ETag != nil {
		out["etag"] = strings.Trim(*put.ETag, "\"")
	}

	return out, nil
}

func (w *s3Writer) ReceiveEvent(e Event) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(&e); err != nil {
		return nil, err
	}
	identifier := ""
	if len(e.Identifier) > 0 {
		identifier = fmt.Sprintf("%s_", e.Identifier)
	}
	identifier = fmt.Sprintf("%s%d", identifier, e.Timestamp.UnixNano())

	return w.uploadToS3(e.ctx, identifier, &buf, e.Timestamp)
}

func (w *s3Writer) Close() {}

func getDateDirectory(now time.Time) string {
	yearNum, monthName, dayNum := now.Date()

	year := strconv.Itoa(yearNum)

	monthPossiblyOneDigit := int(monthName)
	month := fmt.Sprintf("%02d", monthPossiblyOneDigit)

	day := fmt.Sprintf("%02d", dayNum)

	hourPossiblyOneDigit := now.Hour()
	hour := fmt.Sprintf("%02d", hourPossiblyOneDigit)

	return fmt.Sprintf("%s/%s/%s/%s/", year, month, day, hour)
}
