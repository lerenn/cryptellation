package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

type HealthSuite struct {
	suite.Suite
	health *Health
}

func (suite *HealthSuite) SetupTest() {
	// Set env for health
	os.Setenv(HealthPortEnvVar, "1234")

	// Create health
	h, err := NewHealth(context.Background())
	suite.Require().NoError(err)
	suite.health = h
}

func (suite *HealthSuite) TestLiveness() {
	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()
	livenessFunc := suite.health.liveness()

	livenessFunc(w, req)
	suite.Require().Equal(http.StatusOK, w.Result().StatusCode)
}

func (suite *HealthSuite) TestReadiness() {
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	readinessFunc := suite.health.readiness()

	suite.health.Ready(false)
	w := httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusServiceUnavailable, w.Result().StatusCode)

	suite.health.Ready(true)
	w = httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusOK, w.Result().StatusCode)

	suite.health.Ready(false)
	w = httptest.NewRecorder()
	readinessFunc(w, req)
	suite.Require().Equal(http.StatusServiceUnavailable, w.Result().StatusCode)
}
