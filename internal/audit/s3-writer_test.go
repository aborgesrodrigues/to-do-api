package audit

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDirectoryCreation(t *testing.T) {
	now := time.Now()
	actual := getDateDirectory(now)

	if !strings.Contains(actual, strconv.Itoa(now.Year())) {
		t.Errorf("Expected Year: %d in directory: %s", now.Year(), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(int(now.Month()))) {
		t.Errorf("Expected Month: %d in directory: %s", int(now.Month()), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(now.Day())) {
		t.Errorf("Expected Day: %d in directory: %s", now.Day(), actual)
	}

	if !strings.Contains(actual, strconv.Itoa(now.Hour())) {
		t.Errorf("Expected Hour: %d in directory: %s", now.Hour(), actual)
	}
}

func ExampleNewS3Writer() {
	s3Writer, err := NewS3Writer(S3Config{
		// nil or omitted to write to AWS S3. Might set "http://localhost:4566",
		// for example, to write to localstack for local development purposes.
		Endpoint:  nil,
		Region:    "us-east-1",
		Bucket:    "my-audit-logs",
		Directory: "my-service-name",
	})
	if err != nil {
		panic(err)
	}
	logger, err := NewLogger(Config{}, s3Writer)
	defer logger.Close()

	logger.Write(context.Background(), "test", Metadata{Name: "foo", Value: "bar"})
}
