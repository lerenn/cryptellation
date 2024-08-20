package test

import (
	"context"
	"testing"

	"cryptellation/pkg/config"

	client "cryptellation/svc/ticks/clients/go"
	"cryptellation/svc/ticks/clients/go/nats"

	"github.com/stretchr/testify/suite"
)

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

type EndToEndSuite struct {
	suite.Suite
	client client.Client
}

func (suite *EndToEndSuite) SetupSuite() {
	// Get config
	cfg := config.LoadNATS()
	suite.Require().NoError(cfg.Validate())

	// Init client
	client, err := nats.New(cfg)
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *EndToEndSuite) TearDownSuite() {
	suite.client.Close(context.Background())
}
