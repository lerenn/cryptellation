package ticks

import (
	"testing"

	"github.com/lerenn/cryptellation/internal/components/ticks/ports/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestTicksSuite(t *testing.T) {
	suite.Run(t, new(TicksSuite))
}

type TicksSuite struct {
	events.EventsClientSuite
}

func (suite *TicksSuite) SetupTest() {
	ps, err := New(config.LoadNATS())
	suite.Require().NoError(err)
	suite.Client = ps
}
