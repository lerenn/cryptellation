package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/backtests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/backtests/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/client"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/order"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/status"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	"github.com/stretchr/testify/suite"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	vdb       vdb.Port
	client    client.Client
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tmpEnvVar("CRYPTELLATION_BACKTESTS_GRPC_URL", ":9004")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_BACKTESTS_GRPC_URL")
	grpcController := grpc.New(a)
	suite.NoError(grpcController.RunOnAddr(rpcUrl))

	ok := waitForPort(rpcUrl)
	if !ok {
		log.Println("Timed out waiting for trainer gRPC to come up")
	}

	client, closeClient, err := client.New()
	suite.Require().NoError(err)
	suite.client = client

	suite.closeTest = func() error {
		err := closeClient()
		go grpcController.Stop() // TODO: remove goroutine
		closeApplication()
		return err
	}

	vdb, err := redis.New()
	suite.Require().NoError(err)
	suite.vdb = vdb
}

func (suite *ServiceSuite) TearDownSuite() {
	err := suite.closeTest()
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) TestCreateBacktest() {
	accounts := map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{
				"DAI": 1000,
			},
		},
	}

	id, err := suite.client.CreateBacktest(context.Background(), time.Unix(0, 0), time.Unix(120, 0), accounts)
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(id))
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(0, 0), recvBT.StartTime, time.Millisecond)
	suite.Require().WithinDuration(time.Unix(120, 0), recvBT.EndTime, time.Millisecond)
	suite.Require().Len(recvBT.Accounts, 1)
	suite.Require().Len(recvBT.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(1000.0, recvBT.Accounts["exchange"].Balances["DAI"])
}

func (suite *ServiceSuite) TestBacktestSubscribeToEvents() {
	accounts := map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{
				"DAI": 1000,
			},
		},
	}

	id, err := suite.client.CreateBacktest(context.Background(), time.Unix(0, 0), time.Unix(120, 0), accounts)
	suite.Require().NoError(err)

	err = suite.client.SubscribeToBacktestEvents(context.Background(), id, "exchange", "ETH-DAI")
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(id))
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.TickSubscribers, 1)
	suite.Require().Equal("exchange", recvBT.TickSubscribers[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", recvBT.TickSubscribers[0].PairSymbol)
}

func (suite *ServiceSuite) TestBacktestListenEvents() {
	accounts := map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{
				"DAI": 1000,
			},
		},
	}

	id, err := suite.client.CreateBacktest(context.Background(), time.Unix(0, 0), time.Unix(120, 0), accounts)
	suite.Require().NoError(err)

	err = suite.client.SubscribeToBacktestEvents(context.Background(), id, "exchange", "ETH-DAI")
	suite.Require().NoError(err)

	ch, err := suite.client.ListenBacktest(uint(id))
	suite.Require().NoError(err)

	// First candlestick (high)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 2, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// First candlestick (low)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 0.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// First candlestick (close)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// Second candlestick (open)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (high)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 2, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (low)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 0.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (close)
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// End of backtest
	suite.advance(id)
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(120, 0), status.Status{Finished: true})
}

func (suite *ServiceSuite) advance(id uint64) {
	err := suite.client.AdvanceBacktest(context.Background(), id)
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) checkEvent(ch <-chan event.Event, evtType event.Type, t time.Time, content interface{}) {
	evt, ok := <-ch
	suite.Require().True(ok)
	suite.Require().Equal(evtType, evt.Type)
	suite.Require().Equal(t.UTC(), evt.Time)
	suite.Require().Equal(content, evt.Content)
}

func (suite *ServiceSuite) passEvent(ch <-chan event.Event, evtType event.Type) {
	evt, ok := <-ch
	suite.Require().True(ok)
	suite.Require().Equal(evtType, evt.Type)
}

func (suite *ServiceSuite) TestBacktestOrders() {
	accounts := map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{
				"DAI": 1000,
			},
		},
	}

	id, err := suite.client.CreateBacktest(context.Background(), time.Unix(0, 0), time.Unix(600, 0), accounts)
	suite.Require().NoError(err)

	err = suite.client.SubscribeToBacktestEvents(context.Background(), id, "exchange", "ETH-DAI")
	suite.Require().NoError(err)

	err = suite.client.CreateBacktestOrder(context.Background(), id, order.Order{
		Type:         order.TypeIsMarket,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         order.SideIsBuy,
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accounts, err = suite.client.BacktestAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(999.0, accounts["exchange"].Balances["DAI"])
	suite.Require().Equal(1.0, accounts["exchange"].Balances["ETH"])

	ch, err := suite.client.ListenBacktest(uint(id))
	suite.Require().NoError(err)
	for i := 0; i < 5; i++ {
		suite.advance(id)
		suite.passEvent(ch, event.TypeIsTick)
		suite.passEvent(ch, event.TypeIsStatus)
	}

	err = suite.client.CreateBacktestOrder(context.Background(), id, order.Order{
		Type:         order.TypeIsMarket,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         order.SideIsSell,
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accounts, err = suite.client.BacktestAccounts(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Equal(1001.0, accounts["exchange"].Balances["DAI"])
	suite.Require().Equal(0.0, accounts["exchange"].Balances["ETH"])

	orders, err := suite.client.BacktestOrders(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Len(orders, 2)

	suite.Require().WithinDuration(time.Unix(0, 0), *orders[0].ExecutionTime, time.Second)
	suite.Require().Equal(order.TypeIsMarket, orders[0].Type)
	suite.Require().Equal("exchange", orders[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", orders[0].PairSymbol)
	suite.Require().Equal(order.SideIsBuy, orders[0].Side)
	suite.Require().Equal(1.0, orders[0].Quantity)
	suite.Require().Equal(1.0, orders[0].Price)

	suite.Require().WithinDuration(time.Unix(60, 0), *orders[1].ExecutionTime, time.Second)
	suite.Require().Equal(order.TypeIsMarket, orders[1].Type)
	suite.Require().Equal("exchange", orders[1].ExchangeName)
	suite.Require().Equal("ETH-DAI", orders[1].PairSymbol)
	suite.Require().Equal(order.SideIsSell, orders[1].Side)
	suite.Require().Equal(1.0, orders[1].Quantity)
	suite.Require().Equal(2.0, orders[1].Price)
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}

func waitForPort(address string) bool {
	waitChan := make(chan struct{})

	go func() {
		for {
			conn, err := net.DialTimeout("tcp", address, time.Second)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if conn != nil {
				waitChan <- struct{}{}
				return
			}
		}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-waitChan:
		return true
	case <-timeout:
		return false
	}
}
