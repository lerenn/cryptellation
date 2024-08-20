package backtests

import (
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"

	"github.com/lerenn/cryptellation/client/internal/app/ports/events"

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
