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
	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
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
	client    *client.GrpcClient
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
	req := proto.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: map[string]*proto.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().WithinDuration(time.Unix(0, 0), recvBT.StartTime, time.Millisecond)
	suite.Require().WithinDuration(time.Unix(120, 0), recvBT.EndTime, time.Millisecond)
	suite.Require().Len(recvBT.Accounts, 1)
	suite.Require().Len(recvBT.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(1000.0, recvBT.Accounts["exchange"].Balances["DAI"])
}

func (suite *ServiceSuite) TestBacktestSubscribeToEvents() {
	req := proto.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: map[string]*proto.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &proto.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadBacktest(context.Background(), uint(resp.Id))
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.TickSubscribers, 1)
	suite.Require().Equal("exchange", recvBT.TickSubscribers[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", recvBT.TickSubscribers[0].PairSymbol)
}

func (suite *ServiceSuite) TestBacktestListenEvents() {
	req := proto.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(120, 0).Format(time.RFC3339),
		Accounts: map[string]*proto.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &proto.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	ch, err := suite.client.ListenBacktest(uint(resp.Id))
	suite.Require().NoError(err)

	// First candlestick (high)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 2, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// First candlestick (low)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 0.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// First candlestick (close)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(0, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(0, 0), status.Status{Finished: false})

	// Second candlestick (open)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (high)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 2, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (low)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 0.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// Second candlestick (close)
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsTick, time.Unix(60, 0), tick.Tick{PairSymbol: "ETH-DAI", Price: 1.5, Exchange: "exchange"})
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(60, 0), status.Status{Finished: false})

	// End of backtest
	suite.advance(resp.Id)
	suite.checkEvent(ch, event.TypeIsStatus, time.Unix(120, 0), status.Status{Finished: true})
}

func (suite *ServiceSuite) advance(id uint64) {
	_, err := suite.client.AdvanceBacktest(context.Background(), &proto.AdvanceBacktestRequest{
		Id: id,
	})
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
	req := proto.CreateBacktestRequest{
		StartTime: time.Unix(0, 0).Format(time.RFC3339),
		EndTime:   time.Unix(600, 0).Format(time.RFC3339),
		Accounts: map[string]*proto.Account{
			"exchange": {
				Assets: map[string]float32{
					"DAI": 1000,
				},
			},
		},
	}

	resp, err := suite.client.CreateBacktest(context.Background(), &req)
	suite.Require().NoError(err)

	_, err = suite.client.SubscribeToBacktestEvents(context.Background(), &proto.SubscribeToBacktestEventsRequest{
		Id:           resp.Id,
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
	})
	suite.Require().NoError(err)

	_, err = suite.client.CreateBacktestOrder(context.Background(), &proto.CreateBacktestOrderRequest{
		BacktestId:   resp.Id,
		Type:         "market",
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         "buy",
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accountsResp, err := suite.client.Accounts(context.Background(), &proto.AccountsRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(float32(999), accountsResp.Accounts["exchange"].Assets["DAI"])
	suite.Require().Equal(float32(1), accountsResp.Accounts["exchange"].Assets["ETH"])

	ch, err := suite.client.ListenBacktest(uint(resp.Id))
	suite.Require().NoError(err)
	for i := 0; i < 5; i++ {
		suite.advance(resp.Id)
		suite.passEvent(ch, event.TypeIsTick)
		suite.passEvent(ch, event.TypeIsStatus)
	}

	_, err = suite.client.CreateBacktestOrder(context.Background(), &proto.CreateBacktestOrderRequest{
		BacktestId:   resp.Id,
		Type:         "market",
		ExchangeName: "exchange",
		PairSymbol:   "ETH-DAI",
		Side:         "sell",
		Quantity:     1,
	})
	suite.Require().NoError(err)

	accountsResp, err = suite.client.Accounts(context.Background(), &proto.AccountsRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(float32(1001), accountsResp.Accounts["exchange"].Assets["DAI"])
	suite.Require().Equal(float32(0), accountsResp.Accounts["exchange"].Assets["ETH"])

	ordersResp, err := suite.client.Orders(context.Background(), &proto.OrdersRequest{
		BacktestId: resp.Id,
	})
	suite.Require().NoError(err)
	suite.Require().Len(ordersResp.Orders, 2)

	suite.Require().Equal("1970-01-01T00:00:00Z", ordersResp.Orders[0].Time)
	suite.Require().Equal("market", ordersResp.Orders[0].Type)
	suite.Require().Equal("exchange", ordersResp.Orders[0].ExchangeName)
	suite.Require().Equal("ETH-DAI", ordersResp.Orders[0].PairSymbol)
	suite.Require().Equal("buy", ordersResp.Orders[0].Side)
	suite.Require().Equal(float32(1), ordersResp.Orders[0].Quantity)
	suite.Require().Equal(float32(1), ordersResp.Orders[0].Price)

	suite.Require().Equal("1970-01-01T00:01:00Z", ordersResp.Orders[1].Time)
	suite.Require().Equal("market", ordersResp.Orders[1].Type)
	suite.Require().Equal("exchange", ordersResp.Orders[1].ExchangeName)
	suite.Require().Equal("ETH-DAI", ordersResp.Orders[1].PairSymbol)
	suite.Require().Equal("sell", ordersResp.Orders[1].Side)
	suite.Require().Equal(float32(1), ordersResp.Orders[1].Quantity)
	suite.Require().Equal(float32(2), ordersResp.Orders[1].Price)
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
