package nats

import (
	"testing"

	"github.com/lerenn/cryptellation/internal/core/backtests/ports/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestNATSSuite(t *testing.T) {
	suite.Run(t, new(NATSSuite))
}

type NATSSuite struct {
	events.EventsClientSuite
}

func (suite *NATSSuite) SetupTest() {
	ps, err := New(config.LoadNATSConfigFromEnv())
	suite.Require().NoError(err)
	suite.Client = ps
}
