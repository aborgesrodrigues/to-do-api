// Package audit provides an interface for audit logging.
// The audit logger is provided an EventWriter at creation to control what happens to audit events as they arrive.
// Two EventWriters are provided:
//
//	zapWriter writes events to a zap.Logger.
//	s3Writer writes events to an S3 bucket.
//
// A custom EventWriter can be provided by the caller if desired.
// It is the user's responsibility to call Close() on the audit logger instance before shutting down their application to ensure
// writes to the audit logger are flushed and final writes are performed.
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Metadata represents additional fields to include with the audit event.
type Metadata struct {
	Name  string
	Value interface{}
}

// MarshalJSON fulfills the json.Marshaler interface. It is needed to control Metadata JSON output.
func (m *Metadata) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(m.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{"%s":%s}`, m.Name, b)), nil
}

// UnmarshalJSON fulfills the json.Unmarshaler interface. It is needed to control Metadata JSON input.
func (m *Metadata) UnmarshalJSON(b []byte) error {
	var mi map[string]interface{}
	err := json.Unmarshal(b, &mi)
	if err != nil {
		return err
	}
	if len(mi) != 1 {
		return fmt.Errorf("marshal of metadata received unexpected value")
	}
	for k, v := range mi {
		m.Name = k
		m.Value = v
	}
	return nil
}

// Event represents a single audit event.
type Event struct {
	ctx        context.Context
	Identifier string
	Timestamp  time.Time
	Metadata   []Metadata
}

const (
	identifier        = "identifier"
	eventTimestampKey = "event_timestamp"
)

var reservedKeys = [...]string{
	eventTimestampKey,
	identifier,
}

func checkMetadataCollision(key string) error {
	for _, s := range reservedKeys {
		if key == s {
			return fmt.Errorf("metadata attempts to overwrite reserved key: %s", key)
		}
	}
	return nil
}

// MarshalJSON fulfills the json.Marshaler interface. It is needed to control timestamp formatting and write Metadata to the top level JSON object.
func (e *Event) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m[identifier] = e.Identifier
	m[eventTimestampKey] = e.Timestamp.Format(time.RFC3339Nano)
	for _, md := range e.Metadata {
		if err := checkMetadataCollision(md.Name); err != nil {
			return nil, err
		}
		m[md.Name] = md.Value
	}
	return json.Marshal(&m)
}

// EventWriter controls what happens to an audit event once it is received on the events channel.
type EventWriter interface {
	// ReceiveEvent is run each time an audit event is received on the events channel.
	ReceiveEvent(Event) (map[string]interface{}, error)
	// Close is called by the audit logger's Close() method at shutdown to give writers the
	// opportunity to do writer-specific cleanup (e.g., close file, close connection, etc.).
	Close()
}

// Config is a struct containing configuration options for an audit logger.
type Config struct {
	// Defaults to 100.
	BufferSize int
}

// Logger is a struct representing an audit logger.
type Logger struct {
	events      chan Event
	eventWriter EventWriter
	config      Config // copy of the configuration the audit logger was created with
	wg          *sync.WaitGroup

	// Called when an error is encountered while writing audit log to S3.
	onError func(context.Context, error)
	// Called when an audit log is written to S3. Provides object key.
	onWrite func(context.Context, map[string]interface{})
}

func processAuditEvents(logger *Logger) {
	for msg := range logger.events {
		if logger.eventWriter != nil {
			if key, err := logger.eventWriter.ReceiveEvent(msg); err != nil && logger.onError != nil {
				logger.onError(msg.ctx, err)
			} else if logger.onWrite != nil {
				logger.onWrite(msg.ctx, key)
			}
		}
	}
	logger.wg.Done()
}

// NewLogger creates and returns a new audit logger instance.
func NewLogger(config Config, writer EventWriter) (*Logger, error) {
	if config.BufferSize < 1 {
		config.BufferSize = 100
	}
	logger := &Logger{
		events:      make(chan Event, config.BufferSize),
		eventWriter: writer,
		config:      config,
		wg:          &sync.WaitGroup{},
	}

	for range 4 {
		logger.wg.Add(1)
		go processAuditEvents(logger)
	}
	return logger, nil
}

// Close closes the audit logger.
func (logger *Logger) Close() {
	// close stops intake on the channel and will stop workers (goroutines) when channel has flushed
	close(logger.events)

	// waits until channel is flushed or ten seconds have gone by
	waitTimeout(logger.wg, time.Second*10)

	if logger.eventWriter != nil {
		logger.eventWriter.Close()
	}
}

// Write writes an audit log.
func (logger *Logger) Write(ctx context.Context, id string, metadata ...Metadata) {
	logger.events <- Event{
		ctx:        ctx,
		Timestamp:  time.Now(),
		Identifier: id,
		Metadata:   metadata,
	}
}

// SetOnWrite has logger call f when an audit log is written.
func (logger *Logger) SetOnWrite(f func(ctx context.Context, out map[string]interface{})) {
	logger.onWrite = f
}

// SetOnError has logger call f when an error occurs when writing an audit log.
func (logger *Logger) SetOnError(f func(ctx context.Context, err error)) {
	logger.onError = f
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
