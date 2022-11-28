package db

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func TestMain(m *testing.M) {

	logger = zap.NewNop()

	os.Exit(m.Run())
}
