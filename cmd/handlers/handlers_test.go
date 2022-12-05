package handlers

import (
	"testing"

	"github.com/aborgesrodrigues/to-do-api/internal/audit"
	"github.com/aborgesrodrigues/to-do-api/internal/logging"
	mock_service "github.com/aborgesrodrigues/to-do-api/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type handlerTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller // Controller used to create the mock.
	handler *Handler
}

func (hdl *handlerTestSuite) SetupSuite() {
	logger := zap.NewNop()
	auditWriter := audit.NewZapWriter(zaptest.NewLogger(hdl.T()))
	auditLogger, err := logging.NewHTTPAuditLogger(logging.HTTPAuditLogOptions{
		Writer: auditWriter,
	})
	hdl.Assert().NoError(err)

	defer auditLogger.Close()

	hdl.handler = New(logger, auditLogger)
}

func (hdl *handlerTestSuite) SetupTest() {
	hdl.ctrl = gomock.NewController(hdl.Suite.T())
	svcInterface := mock_service.NewMockSVCInterface(hdl.ctrl)
	hdl.handler.svc = svcInterface
}

func (hdl *handlerTestSuite) TearDownTest() {
	hdl.ctrl.Finish()
}

func (hdl *handlerTestSuite) getService() *mock_service.MockSVCInterfaceMockRecorder {
	return hdl.handler.svc.(*mock_service.MockSVCInterface).EXPECT()
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
