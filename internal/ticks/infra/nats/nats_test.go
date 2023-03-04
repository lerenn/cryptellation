package nats

import (
	"testing"

	"github.com/digital-feather/cryptellation/internal/ticks/app/ports/pubsub/tests"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/stretchr/testify/suite"
)

func TestNATSSuite(t *testing.T) {
	suite.Run(t, new(NATSSuite))
}

type NATSSuite struct {
	tests.PubSubClientSuite
}

func (suite *NATSSuite) SetupTest() {
	ps, err := New(config.LoadNATSConfigFromEnv())
	suite.Require().NoError(err)
	suite.Client = ps
}
