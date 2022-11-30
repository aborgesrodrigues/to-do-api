package db

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type dbTestSuite struct {
	suite.Suite
	db   *DB
	mock sqlmock.Sqlmock
}

func (d *dbTestSuite) SetupSuite() {
	logger := zap.NewNop()

	cfg := Config{
		Logger: logger,
	}

	dbMock, mock, err := sqlmock.New()
	d.Assert().NoError(err)

	d.db, err = New(cfg)
	d.Assert().NoError(err)

	d.db.db = dbMock
	d.mock = mock
}

func (d *dbTestSuite) TearDownSuite() {
	d.db.db.Close()
}

func (d *dbTestSuite) SetupTest() {

}

func (d *dbTestSuite) TearDownTest() {

}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbTestSuite))
}
