package test

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

func (suite *EndToEndSuite) TestListExchanges() {
	// WHEN requesting the exchanges list

	res, err := suite.client.ListExchanges(context.Background(), api.ListExchangesWorkflowParams{})

	// THEN the request is successful

	suite.Require().NoError(err)

	// AND the response contains the proper exchanges

	suite.Require().Equal([]string{"binance"}, res.List)
}

func (suite *EndToEndSuite) TestGetExchange() {
	// WHEN requesting an exchange

	res, err := suite.client.GetExchange(context.Background(), api.GetExchangeWorkflowParams{
		Name: "binance",
	})

	// THEN the request is successful

	suite.Require().NoError(err)

	// AND the response contains the proper exchange

	suite.Require().Equal("binance", res.Exchange.Name)
	suite.Require().NotEmpty(res.Exchange.Pairs)
	suite.Require().NotEmpty(res.Exchange.Periods)
}
