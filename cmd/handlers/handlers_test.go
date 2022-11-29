package handlers

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	hdl    *Handler
)

func TestMain(m *testing.M) {

	logger = zap.NewNop()

	hdl = New(logger)

	os.Exit(m.Run())
}
