package main

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	hdl    *handler
)

func TestMain(m *testing.M) {

	logger = zap.NewNop()

	hdl = newHandler()

	os.Exit(m.Run())
}
