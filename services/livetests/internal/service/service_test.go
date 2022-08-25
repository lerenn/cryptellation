package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb"
	"github.com/digital-feather/cryptellation/services/livetests/internal/adapters/vdb/redis"
	"github.com/digital-feather/cryptellation/services/livetests/internal/controllers/grpc"
	"github.com/digital-feather/cryptellation/services/livetests/pkg/client"
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
	defer tmpEnvVar("CRYPTELLATION_LIVETESTS_GRPC_URL", ":9006")()

	a, closeApplication, err := NewMockedApplication()
	suite.Require().NoError(err)

	rpcUrl := os.Getenv("CRYPTELLATION_LIVETESTS_GRPC_URL")
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

func (suite *ServiceSuite) TestCreateLivetest() {
	id, err := suite.client.CreateLivetest(context.Background(), map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{
				"DAI": 1000,
			},
		},
	})
	suite.Require().NoError(err)

	recvBT, err := suite.vdb.ReadLivetest(context.Background(), id)
	suite.Require().NoError(err)
	suite.Require().Len(recvBT.Accounts, 1)
	suite.Require().Len(recvBT.Accounts["exchange"].Balances, 1)
	suite.Require().Equal(1000.0, recvBT.Accounts["exchange"].Balances["DAI"])
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
