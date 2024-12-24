package test

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/v1/clients/go/direct"
	"github.com/stretchr/testify/suite"
)

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

type EndToEndSuite struct {
	suite.Suite
	client direct.Client
}

func (suite *EndToEndSuite) SetupSuite() {
	client, err := direct.NewClient()
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *EndToEndSuite) TearDownSuite() {
	suite.client.Close(context.Background())
}
