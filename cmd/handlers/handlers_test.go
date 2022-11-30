package handlers

import (
	"testing"

	mock_service "github.com/aborgesrodrigues/to-do-api/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type handlerTestSuite struct {
	suite.Suite
	ctrl    *gomock.Controller // Controller used to create the mock.
	handler *Handler
}

func (hdl *handlerTestSuite) SetupSuite() {
	logger := zap.NewNop()

	hdl.handler = New(logger)
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
