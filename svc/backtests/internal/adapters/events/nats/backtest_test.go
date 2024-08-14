package backtests

import (
	"testing"

	"cryptellation/internal/config"

	"cryptellation/svc/backtests/internal/app/ports/events"

	"github.com/stretchr/testify/suite"
)

func TestBacktestsSuite(t *testing.T) {
	suite.Run(t, new(BacktestsSuite))
}

type BacktestsSuite struct {
	events.EventsClientSuite
}

func (suite *BacktestsSuite) SetupTest() {
	ps, err := New(config.LoadNATS())
	suite.Require().NoError(err)
	suite.Client = ps
}
