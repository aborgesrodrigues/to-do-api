package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aborgesrodrigues/to-do-api/cmd/handlers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type testSuite struct {
	suite.Suite
	handler *handlers.Handler
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &testSuite{})
}

func (s *testSuite) SetupSuite() {
	viper.AutomaticEnv()
	viper.SetConfigFile("../../.env")
	viper.ReadInConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Error creating logger")
	}

	hdl := handlers.New(logger)
	s.handler = hdl

	logger.Info("Server listening.", zap.String("addr", "8080"))
	go func() {
		if err := http.ListenAndServe(":8080", getRouter(hdl)); err != nil {
			logger.Error(err.Error())
		}
	}()
}

func (s *testSuite) TearDownSuite() {
	// p, _ := os.FindProcess(syscall.Getpid())
	// p.Signal(syscall.SIGINT)
}

func (s *testSuite) SetupTest() {

}

func (s *testSuite) TearDownTest() {

}

func (s *testSuite) call(method, path string, payload, response interface{}) {
	var buf bytes.Buffer
	if payload != nil {
		err := json.NewEncoder(&buf).Encode(payload)
		s.Assert().NoError(err)
	}

	req, err := http.NewRequest(method, path, nil)
	s.Assert().NoError(err)

	res, err := http.DefaultClient.Do(req)
	s.Assert().NoError(err)

	if response != nil {
		err = json.NewDecoder(res.Body).Decode(&response)
		s.Assert().NoError(err)
	}
}
