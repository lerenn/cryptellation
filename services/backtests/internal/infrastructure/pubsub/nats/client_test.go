package nats

import (
	"testing"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/pubsub/tests"
	"github.com/stretchr/testify/suite"
)

func TestNATSClientSuite(t *testing.T) {
	suite.Run(t, new(NATSClientSuite))
}

type NATSClientSuite struct {
	tests.PubSubClientSuite
}

func (suite *NATSClientSuite) SetupTest() {
	client, err := New()
	suite.Require().NoError(err)
	suite.Client = client
}
