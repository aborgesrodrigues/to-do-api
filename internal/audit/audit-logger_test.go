package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventJSONMarshal(t *testing.T) {
	timestamp := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	for name, tt := range map[string]struct {
		event       Event
		expected    string
		expectedErr string
	}{
		"happy": {
			event: Event{
				Timestamp:  timestamp,
				Identifier: "12345",
				Metadata:   []Metadata{{Name: "foo", Value: "test"}},
			},
			expected: `{"event_timestamp":"1970-01-01T00:00:00Z","foo":"test","identifier":"12345"}`,
		},
		"metadata collision": {
			event: Event{
				Timestamp: timestamp,
				Metadata: []Metadata{
					{Name: "identifier", Value: "test"},
					{Name: "event_timestamp", Value: "test"},
				},
			},
			expectedErr: "json: error calling MarshalJSON for type *audit.Event: metadata attempts to overwrite reserved key: identifier",
		},
	} {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(&tt.event)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
			assert.Equal(t, tt.expected, string(b))
		})
	}
}

func TestEventJSONUnmarshal(t *testing.T) {

	// Create type with nested types to test nested marshal/unmarshal
	type Person struct {
		Name        string
		Address     string
		Friends     []Person
		Connections map[string]string
	}

	timestamp := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	for name, tt := range map[string]struct {
		event interface{}
	}{
		"happy": {
			event: Event{
				Timestamp:  timestamp,
				Identifier: "12345",
				Metadata: []Metadata{
					{
						Name: "test",
						Value: Person{
							Address: "123 Fake",
							Connections: map[string]string{
								"jake":  "fake",
								"other": "data",
							},
							Friends: []Person{
								{
									Name:    "Tom",
									Address: "456 Fake",
								},
								{
									Name:    "James",
									Address: "333 Fake",
								},
							},
							Name: "Geroge",
						},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// marshal orig to get bytes
			origBytes, err := json.Marshal(&tt.event)
			assert.NoError(t, err)
			origReader := bytes.NewReader(origBytes)
			decoder := json.NewDecoder(origReader)

			// unmarshal into event struct
			var origUnmarshalEvent Event
			err = decoder.Decode(&origUnmarshalEvent)
			assert.NoError(t, err)

			// marshal the event that has done unmarshal-ing
			newBytes, err := json.Marshal(origUnmarshalEvent)
			assert.NoError(t, err)
			newReader := bytes.NewReader(newBytes)
			decoder = json.NewDecoder(newReader)

			// unmarshal the bytes again
			var newUnmarshalEvent Event
			err = decoder.Decode(&newUnmarshalEvent)
			if err != nil {
				panic(err)
			}

			// marshal -> unmarshal two times equality proves uniformity
			fmt.Println(reflect.DeepEqual(origUnmarshalEvent, newUnmarshalEvent))
		})
	}
}

func TestMetadataJSONUnmarshal(t *testing.T) {
	for name, tt := range map[string]struct {
		metadata    interface{}
		expectedErr bool
	}{
		"happy": {
			metadata: map[string]interface{}{"test": "test"},
		},
		"unmarshalErr": {
			metadata:    "test",
			expectedErr: true,
		},
		"tooManyFields": {
			metadata:    map[string]interface{}{"test": "test", "another": "another"},
			expectedErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			// marshalBytes is marshaled val
			marshalBytes, err := json.Marshal(tt.metadata)
			assert.NoError(t, err)
			reader := bytes.NewReader(marshalBytes)

			// attempt to unmarshal
			decoder := json.NewDecoder(reader)
			var origUnmarshalEvent Metadata
			err = decoder.Decode(&origUnmarshalEvent)
			assert.Equal(t, tt.expectedErr, err != nil)
		})
	}
}

type mockWriter struct {
	mock.Mock
}

func (w *mockWriter) OnReceiveEvent(arguments ...interface{}) *mock.Call {
	return w.On("ReceiveEvent", arguments...)
}

func (w *mockWriter) OnClose(arguments ...interface{}) *mock.Call {
	return w.On("Close", arguments...)
}

func (w *mockWriter) ReceiveEvent(e Event) (map[string]interface{}, error) {
	if e.Timestamp.IsZero() {
		panic("timestamp shouldn't be zero")
	}
	// HAAACK
	e.Timestamp = time.Time{}
	args := w.Called(e)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (w *mockWriter) Close() {
	w.Called()
}

func TestErrorHandler(t *testing.T) {
	e := Event{
		ctx:        context.Background(),
		Identifier: "id",
		Metadata:   []Metadata{{Name: "foo", Value: "bar"}},
	}

	writer := mockWriter{}
	writer.
		OnReceiveEvent(e).
		Return(map[string]interface{}{}, errors.New("test"))
	writer.OnClose()

	logger, err := NewLogger(Config{
		BufferSize: 1,
	}, &writer)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	onErrorCalled := false
	logger.SetOnError(func(ctx context.Context, err error) {
		onErrorCalled = true
		assert.EqualError(t, err, "test")
	})

	logger.Write(context.Background(), "id", Metadata{Name: "foo", Value: "bar"})
	logger.Close()

	writer.AssertCalled(t, "ReceiveEvent", e)
	assert.True(t, onErrorCalled, "error handler should be called")
}

func TestWriteHandler(t *testing.T) {
	e := Event{
		ctx:        context.Background(),
		Identifier: "id",
		Metadata:   []Metadata{{Name: "foo", Value: "bar"}},
	}

	writer := mockWriter{}
	writer.
		OnReceiveEvent(e).
		Return(map[string]interface{}{
			"baz": "qux",
		}, nil)
	writer.OnClose()

	logger, err := NewLogger(Config{
		BufferSize: 1,
	}, &writer)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	onWriteCalled := false
	logger.SetOnWrite(func(ctx context.Context, out map[string]interface{}) {
		onWriteCalled = true
		assert.Equal(t, out, map[string]interface{}{
			"baz": "qux",
		})
	})

	logger.Write(context.Background(), "id", Metadata{Name: "foo", Value: "bar"})
	logger.Close()

	writer.AssertCalled(t, "ReceiveEvent", e)
	assert.True(t, onWriteCalled, "write handler should be called")
}
