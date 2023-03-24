package nats

import (
	"testing"

	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/services/ticks/io/events/tests"
	"github.com/stretchr/testify/suite"
)

func TestNATSSuite(t *testing.T) {
	suite.Run(t, new(NATSSuite))
}

type NATSSuite struct {
	tests.EventsClientSuite
}

func (suite *NATSSuite) SetupTest() {
	ps, err := New(config.LoadNATSConfigFromEnv())
	suite.Require().NoError(err)
	suite.Client = ps
}
