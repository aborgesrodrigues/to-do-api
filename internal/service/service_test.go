package service

import (
	"testing"

	mock_db "github.com/aborgesrodrigues/to-do-api/internal/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type svcTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller // Controller used to create the mock.
	svc  *Service
}

func (s *svcTestSuite) SetupSuite() {
	logger := zap.NewNop()

	svc, err := New(Config{Logger: logger})
	s.Assert().NoError(err)

	s.svc = svc
}

func (s *svcTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.Suite.T())
	dbInterface := mock_db.NewMockDBInterface(s.ctrl)
	s.svc.db = dbInterface
}

func (s *svcTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *svcTestSuite) getDB() *mock_db.MockDBInterfaceMockRecorder {
	return s.svc.db.(*mock_db.MockDBInterface).EXPECT()
}

func TestService(t *testing.T) {
	suite.Run(t, new(svcTestSuite))
}
