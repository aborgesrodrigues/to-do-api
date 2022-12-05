package audit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapWriter(t *testing.T) {
	core, observed := observer.New(zap.DebugLevel)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return core
	})))

	writer := NewZapWriter(logger)

	e := Event{
		ctx:        context.Background(),
		Identifier: "id",
		Timestamp:  time.Time{},
		Metadata: []Metadata{
			{Name: "foo", Value: "bar"},
		},
	}
	out, err := writer.ReceiveEvent(e)
	writer.Close()

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"message": "Audit event."}, out)
	expected := []observer.LoggedEntry{
		{
			Entry: zapcore.Entry{
				Level:   zap.InfoLevel,
				Message: "Audit event.",
			},
			Context: []zapcore.Field{
				{
					Key:    "identifier",
					Type:   zapcore.StringType,
					String: e.Identifier,
				},
				{
					Key:       "metadata",
					Type:      zapcore.ReflectType,
					Interface: e.Metadata,
				},
			},
		},
	}
	assert.Equal(t, expected, observed.AllUntimed())
}
