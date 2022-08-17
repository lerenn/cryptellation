package nats

import (
	"os"
	"testing"

	"github.com/digital-feather/cryptellation/services/ticks/internal/adapters/pubsub/tests"
	"github.com/stretchr/testify/suite"
)

func TestNATSClientSuite(t *testing.T) {
	if os.Getenv("NATS_URL") == "" {
		t.Skip()
	}

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
