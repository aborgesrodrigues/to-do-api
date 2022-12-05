package audit

import "context"

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
