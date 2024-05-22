package domain

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetAccountsSuite(t *testing.T) {
	suite.Run(t, new(GetAccountsSuite))
}

type GetAccountsSuite struct {
	suite.Suite
	operator     app.ForwardTests
	db           *db.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *GetAccountsSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.candlesticks)
}

func (suite *GetAccountsSuite) TestHappyPass() {
	id := uuid.New()

	// Set DB mock expectations
	suite.db.EXPECT().ReadForwardTest(context.Background(), id).
		Return(forwardtest.ForwardTest{
			Accounts: map[string]account.Account{
				"exchange": {
					Balances: map[string]float64{"DAI": 1000},
				},
			},
		}, nil)

	// Execute getting accounts
	accounts, err := suite.operator.GetAccounts(context.Background(), id)
	suite.Require().NoError(err)

	suite.Require().Equal(map[string]account.Account{
		"exchange": {
			Balances: map[string]float64{"DAI": 1000},
		},
	}, accounts)
}
