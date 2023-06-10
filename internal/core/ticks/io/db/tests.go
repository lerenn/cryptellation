package db

import (
	"context"

	"github.com/stretchr/testify/suite"
)

type SymbolListenerSuite struct {
	suite.Suite
	DB Port
}

func (suite *SymbolListenerSuite) TestIncrementDecrementSymbolListener() {
	err := suite.DB.ClearSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	err = suite.DB.ClearSymbolListenerSubscribers(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)

	count, err := suite.DB.IncrementSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.DB.IncrementSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(2), count)

	count, err = suite.DB.IncrementSymbolListenerSubscribers(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.DB.GetSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(2), count)

	count, err = suite.DB.GetSymbolListenerSubscribers(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	count, err = suite.DB.DecrementSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(1), count)

	suite.NoError(suite.DB.ClearAllSymbolListenersCount(context.Background()))

	count, err = suite.DB.GetSymbolListenerSubscribers(context.Background(), "EXCHANGE1", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(0), count)

	count, err = suite.DB.GetSymbolListenerSubscribers(context.Background(), "EXCHANGE2", "PAIR1")
	suite.Require().NoError(err)
	suite.Equal(int64(0), count)
}
