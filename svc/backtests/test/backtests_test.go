package test

import (
	"context"
	"time"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/pkg/utils"

	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (suite *EndToEndSuite) TestBacktestGet() {
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

	// Get backtest
	b, err := suite.client.Get(context.Background(), id)
	suite.Require().NoError(err)

	// Check backtest
	suite.Require().Equal(id, b.ID)
	suite.Require().Equal(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:00:00Z")), b.Parameters.StartTime)
	suite.Require().Equal(utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:02:00Z")), b.Parameters.EndTime)
	suite.Require().Equal([]event.PricesSubscription{{Exchange: "binance", Pair: "BTC-USDT"}}, b.PricesSubscriptions)
}

func (suite *EndToEndSuite) TestBacktestAdvanceWithFullOHLC() {
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
		Mode:      backtest.ModeIsFullOHLC.Opt(),
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

	// Getting price event
	suite.checkPriceEvent(t, tick.Tick{
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

func (suite *EndToEndSuite) TestBacktestAdvanceWithCloseOHLC() {
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
		Mode:      backtest.ModeIsCloseOHLC.Opt(),
	})
	suite.Require().NoError(err)

	// Subscribe to pair
	suite.Require().NoError(suite.client.Subscribe(context.Background(), id, "binance", "BTC-USDT"))

	// Listen to events
	ch, err := suite.client.ListenEvents(context.Background(), id)
	suite.Require().NoError(err)

	// Advance for 12:01 (Close price of candlestick)
	err = suite.client.Advance(context.Background(), id)
	suite.Require().NoError(err)

	// Getting price event
	t := utils.Must(time.Parse(time.RFC3339, "2023-02-26T12:01:00Z"))
	suite.checkPriceEvent(t, tick.Tick{
		Time:     t,
		Pair:     "BTC-USDT",
		Price:    23272.77,
		Exchange: "binance",
	}, <-ch)

	// Getting status event
	suite.checkStatusEvent(t, false, <-ch)
}

func (suite *EndToEndSuite) checkPriceEvent(expectedTime time.Time, expectedTick tick.Tick, receivedEvt event.Event) {
	// Check type and time
	suite.Require().Equal(event.TypeIsPrice, receivedEvt.Type)
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
		Mode:      backtest.ModeIsFullOHLC.Opt(),
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
	err = suite.client.CreateOrder(context.Background(), common.OrderCreationPayload{
		RunID:    id,
		Type:     order.TypeIsMarket,
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Side:     order.SideIsSell,
		Quantity: 1,
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
	err = suite.client.CreateOrder(context.Background(), common.OrderCreationPayload{
		RunID:    id,
		Type:     order.TypeIsMarket,
		Exchange: "binance",
		Pair:     "BTC-USDT",
		Side:     order.SideIsBuy,
		Quantity: 23253.40 / 23269.07, // BTC to buy => USDT Qty / Price
	})
	suite.Require().NoError(err)

	// Check account
	accounts, err = suite.client.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(0.99933, utils.Round(accounts["binance"].Balances["BTC"], 5))
	suite.Require().Equal(0.04, utils.Round(accounts["binance"].Balances["USDT"], 2))

	// List orders
	orders, err := suite.client.ListOrders(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Len(orders, 2)
}

func (suite *EndToEndSuite) TestBacktestList() {
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
		Mode:      backtest.ModeIsFullOHLC.Opt(),
	})
	suite.Require().NoError(err)

	// List backtests
	backtests, err := suite.client.List(context.Background())
	suite.Require().NoError(err)

	found := false
	for _, b := range backtests {
		if b.ID == id {
			found = true
		}
	}

	if !found {
		suite.Fail("Backtest not found")
	}
}
