package test

// func (suite *TicksSuite) TestManualCleanUp() {
// 	err := suite.client.Unregister(context.Background(), client.TicksFilterPayload{
// 		Exchange: "binance",
// 		Pair:   "BTC-USDT",
// 	})
// 	suite.Require().NoError(err)
// }

// TODO: fix listening
// func (suite *TicksSuite) TestListen() {
// 	// Register listener
// 	err := suite.client.Register(context.Background(), client.TicksFilterPayload{
// 		Exchange: "binance",
// 		Pair:   "BTC-USDT",
// 	})
// 	suite.Require().NoError(err)

// 	// Listen to ticks
// 	ch, err := suite.client.Listen(context.Background(), client.TicksFilterPayload{
// 		Exchange: "binance",
// 		Pair:   "BTC-USDT",
// 	})
// 	suite.Require().NoError(err)

// 	// Check that ticks are correct
// 	for i := 0; i < 3; i++ {
// 		t := <-ch
// 		suite.Require().Equal("binance", t.Exchange)
// 		suite.Require().Equal("BTC-USDT", t.Pair)
// 		suite.Require().NotEqual(0, t.Price)
// 		suite.Require().WithinDuration(time.Now(), t.Time, time.Second)
// 	}

// 	// Unregister listener
// 	err = suite.client.Unregister(context.Background(), client.TicksFilterPayload{
// 		Exchange: "binance",
// 		Pair:   "BTC-USDT",
// 	})
// 	suite.Require().NoError(err)
// }
