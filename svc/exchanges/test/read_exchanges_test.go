package test

import "context"

func (suite *EndToEndSuite) TestReadExchanges() {
	// WHEN requesting a exchanges list
	list, err := suite.client.Read(context.Background(), "binance")

	// THEN the request is successful
	suite.Require().NoError(err)

	// AND the response contains the proper exchanges
	suite.Require().Len(list, 1)
	suite.Require().Equal("binance", list[0].Name)

	l := []string{"D1", "D3", "H1", "H12", "H2", "H4", "H6", "H8", "M1", "M15", "M3", "M30", "M5", "W1"}
	for i, s := range l {
		suite.Require().Contains(list[0].PeriodsSymbols, s, i)
	}
}
