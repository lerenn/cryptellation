package domain

import (
	"context"
	"testing"

	candlesticks "cryptellation/svc/candlesticks/clients/go"

	"cryptellation/svc/forwardtests/internal/app"
	"cryptellation/svc/forwardtests/internal/app/ports/db"
	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

type ListSuite struct {
	suite.Suite
	operator     app.ForwardTests
	db           *db.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *ListSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.candlesticks)
}

func (suite *ListSuite) TestHappyPass() {
	id1 := uuid.New()
	id2 := uuid.New()

	// Set DB mock expectations
	suite.db.EXPECT().ListForwardTests(context.Background(), db.ListFilters{}).
		Return([]forwardtest.ForwardTest{
			{ID: id1},
			{ID: id2},
		}, nil)

	// Execute getting forward tests
	forwardTests, err := suite.operator.List(context.Background(), app.ListFilters{})
	suite.Require().NoError(err)

	suite.Require().Equal([]forwardtest.ForwardTest{
		{ID: id1},
		{ID: id2},
	}, forwardTests)
}
