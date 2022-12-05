package audit

import (
	"go.uber.org/zap"
)

type zapWriter struct {
	logger *zap.Logger
}

// NewZapWriter provides a zap-based EventWriter. This is most useful in unit tests where an
// audit.Logger is required. In such a case, creating the zap logger using zaptest is recommended.
func NewZapWriter(logger *zap.Logger) EventWriter {
	return &zapWriter{
		logger: logger,
	}
}

func (s zapWriter) ReceiveEvent(e Event) (map[string]interface{}, error) {
	const message = "Audit event."
	// Not logging event timestamp due to non-determinism and likelihood of this writer being used
	// in unit tests.
	s.logger.Info(message,
		zap.String("identifier", e.Identifier),
		zap.Any("metadata", e.Metadata),
	)
	return map[string]interface{}{"message": message}, nil
}

func (s zapWriter) Close() {
	s.logger.Sync()
}
