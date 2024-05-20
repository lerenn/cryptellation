package test

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/event"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (suite *EndToEndSuite) TestBacktestAdvance() {
	// Create backtest
	id, err := suite.client.Create(context.Background(), client.BacktestCreationPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"BTC": 1,
				},
			},
		},
		StartTime: utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")),
		EndTime:   utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z"))),
	})
	suite.Require().NoError(err)

	// Subscribe to pair
	suite.Require().NoError(suite.client.Subscribe(context.Background(), id, "binance", "BTC-USDT"))

	// Listen to events
	ch, err := suite.client.ListenEvents(context.Background(), id)
	suite.Require().NoError(err)

	// Advance for 12:00 (High price of candlestick)
	t := utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z"))
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23255.54,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:00 (Low price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23248.96,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:00 (Close price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23253.8,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:01 (Open price of candlestick)
	t = utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:01:00Z"))
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23254.26,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:01 (High price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23272.86,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:01 (Low price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23250.65,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:01 (Close price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting tick event
	suite.checkTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23272.77,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)

	// Advance for 12:01 (Open price of candlestick)
	t = utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z"))
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting status event
	suite.checkStatusEvent(t, true, <-ch)
}

func (suite *EndToEndSuite) checkTickEvent(expectedTime time.Time, expectedTick tick.Tick, receivedEvt event.Event) {
	// Check type and time
	suite.Require().Equal(event.TypeIsTick, receivedEvt.Type)
	suite.Require().WithinDuration(expectedTime, receivedEvt.Time, time.Second)

	// Check content
	t, ok := receivedEvt.Content.(tick.Tick)
	suite.Require().True(ok)
	suite.Require().Equal(expectedTick, t)
}

func (suite *EndToEndSuite) checkStatusEvent(expectedTime time.Time, isFinished bool, receivedEvt event.Event) {
	// Check type and time
	suite.Require().Equal(event.TypeIsStatus, receivedEvt.Type)
	suite.Require().WithinDuration(expectedTime, receivedEvt.Time, time.Second)

	// Check content
	t, ok := receivedEvt.Content.(event.Status)
	suite.Require().True(ok)
	suite.Require().Equal(isFinished, t.Finished)
}

func (suite *EndToEndSuite) TestBacktestOrder() {
	// Create backtest
	id, err := suite.client.Create(context.Background(), client.BacktestCreationPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"BTC": 1,
				},
			},
		},
		StartTime: utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")),
		EndTime:   utils.ToReference(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:30:00Z"))),
	})
	suite.Require().NoError(err)

	// Subscribe to pair
	suite.Require().NoError(suite.client.Subscribe(context.Background(), id, "binance", "BTC-USDT"))

	// Listen to events
	ch, err := suite.client.ListenEvents(context.Background(), id)
	suite.Require().NoError(err)

	// Check account
	accounts, err := suite.client.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(1.0, accounts["binance"].Balances["BTC"])

	// Send order
	err = suite.client.CreateOrder(context.Background(), client.OrderCreationPayload{
		BacktestID: id,
		Type:       order.TypeIsMarket,
		Exchange:   "binance",
		Pair:       "BTC-USDT",
		Side:       order.SideIsSell,
		Quantity:   1,
	})
	suite.Require().NoError(err)

	// Check account
	accounts, err = suite.client.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(0.0, accounts["binance"].Balances["BTC"])
	suite.Require().Equal(23253.44, accounts["binance"].Balances["USDT"])

	// Loop over ticks
	for i := 0; i < 4*3; i++ {
		// Advance
		err = suite.client.Advance(context.Background(), id)
		suite.Require().NoError(err)

		// Wait for tick and status to pass
		<-ch
		<-ch
	}

	// Send order
	err = suite.client.CreateOrder(context.Background(), client.OrderCreationPayload{
		BacktestID: id,
		Type:       order.TypeIsMarket,
		Exchange:   "binance",
		Pair:       "BTC-USDT",
		Side:       order.SideIsBuy,
		Quantity:   23253.40 / 23269.07, // BTC to buy => USDT Qty / Price
	})
	suite.Require().NoError(err)

	// Check account
	accounts, err = suite.client.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(0.99933, utils.Round(accounts["binance"].Balances["BTC"], 5))
	suite.Require().Equal(0.04, utils.Round(accounts["binance"].Balances["USDT"], 2))
}
