package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	sqldb "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db/sql"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/client"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
	"github.com/stretchr/testify/suite"
)

const (
	testDatabase = "candlesticks"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	db        db.Port
	client    client.Client
	closeTest func() error
}

func (suite *ServiceSuite) SetupSuite() {
	defer tmpEnvVar("SQLDB_DATABASE", testDatabase)()
	defer tmpEnvVar("CRYPTELLATION_CANDLESTICKS_GRPC_URL", ":9002")()

	a, err := newMockApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_CANDLESTICKS_GRPC_URL")
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
		go grpcController.Stop() // TODO: remove goroutine
		return closeClient()
	}
}

func (suite *ServiceSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", testDatabase)()

	db, err := sqldb.New()
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset())

	suite.db = db
}

func (suite *ServiceSuite) TearDownSuite() {
	err := suite.closeTest()
	suite.Require().NoError(err)
}

func (suite *ServiceSuite) TestGetCandlesticksAllExistWithNoneInDB() {
	// Given a client service
	// Provided before

	// When a request is made
	l, err := suite.client.ReadCandlesticks(context.Background(), client.ReadCandlestickPayload{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        time.Unix(0, 0),
		End:          time.Unix(540, 0),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Equal(10, l.Len())
	i := 0
	_ = l.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().Equal(float64(60*i), cs.Open)
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		i++

		return false, nil
	})
}

func (suite *ServiceSuite) TestGetCandlesticksAllInexistantWithNoneInDB() {
	// Given a client service
	// Provided before

	// When a request is made
	l, err := suite.client.ReadCandlesticks(context.Background(), client.ReadCandlestickPayload{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        time.Unix(60000, 0),
		End:          time.Unix(60600, 0),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Equal(0, l.Len())
}

func (suite *ServiceSuite) TestGetCandlesticksFromDBAndService() {
	// Given a client service
	// Provided before

	// And candlesticks in DB
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	for i := int64(0); i < 10; i++ {
		err := cl.Set(time.Unix(60*i, 0), candlestick.Candlestick{
			Open:  float64(i),
			Close: 4321,
		})
		suite.Require().NoError(err)
	}
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), cl))

	// When a request is made
	l, err := suite.client.ReadCandlesticks(context.Background(), client.ReadCandlestickPayload{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        time.Unix(0, 0),
		End:          time.Unix(1140, 0),
	})

	// Then all candlesticks are retrieved
	suite.Require().NoError(err)
	suite.Require().Equal(20, l.Len())
	i := 0
	_ = l.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		if i < 10 {
			suite.Require().Equal(float64(4321), cs.Close, i)
		} else {
			suite.Require().Equal(float64(1234), cs.Close, i)
		}
		i++
		return false, nil
	})
}

func (suite *ServiceSuite) TestGetCandlesticksFromDBAndServiceWithUncomplete() {
	// Given a client service
	// Provided before

	// And candlesticks in DB
	cl := candlestick.NewList(candlestick.ListID{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
	})
	for i := int64(0); i < 10; i++ {
		cs := candlestick.Candlestick{
			Open:  float64(i),
			Close: 4321,
		}

		if i == 5 {
			cs.Uncomplete = true
		}

		err := cl.Set(time.Unix(60*i, 0), cs)
		suite.Require().NoError(err)
	}
	suite.Require().NoError(suite.db.CreateCandlesticks(context.Background(), cl))

	// When a request is made
	l, err := suite.client.ReadCandlesticks(context.Background(), client.ReadCandlestickPayload{
		ExchangeName: "mock_exchange",
		PairSymbol:   "ETH-USDC",
		Period:       period.M1,
		Start:        time.Unix(0, 0),
		End:          time.Unix(1140, 0),
	})

	// Then all candlesticks are retrieved from mock
	suite.Require().NoError(err)
	suite.Require().Equal(20, l.Len())
	i := 0
	_ = l.Loop(func(t time.Time, cs candlestick.Candlestick) (bool, error) {
		suite.Require().WithinDuration(time.Unix(int64(60*i), 0), t, time.Millisecond)
		suite.Require().Equal(float64(1234), cs.Close, i)
		i++
		return false, nil
	})
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
