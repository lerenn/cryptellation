package test

import (
	"context"
	"testing"

	client "github.com/lerenn/cryptellation/v1/clients/go"
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
	client, err := client.New()
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *EndToEndSuite) TearDownSuite() {
	suite.client.Close(context.Background())
}
