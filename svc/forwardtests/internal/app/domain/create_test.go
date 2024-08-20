package domain

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/models/account"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"

	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestCreationSuite(t *testing.T) {
	suite.Run(t, new(CreationSuite))
}

type CreationSuite struct {
	suite.Suite
	operator     app.ForwardTests
	db           *db.MockPort
	candlesticks *candlesticks.MockClient
}

func (suite *CreationSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.candlesticks = candlesticks.NewMockClient(gomock.NewController(suite.T()))
	suite.operator = New(suite.db, suite.candlesticks)
}

func (suite *CreationSuite) TestHappyPass() {
	ctx := context.Background()
	var appSetID uuid.UUID

	// Set DB mock expectations
	suite.db.EXPECT().CreateForwardTest(ctx, gomock.Any()).
		Do(func(ctx context.Context, ft forwardtest.ForwardTest) {
			appSetID = ft.ID

			suite.Require().Equal(forwardtest.ForwardTest{
				ID: ft.ID,
				Accounts: map[string]account.Account{
					"exchange": {
						Balances: map[string]float64{"DAI": 1000},
					},
				},
			}, ft)
		}).
		Return(nil)

	// Execute creation
	id, err := suite.operator.Create(ctx, forwardtest.NewPayload{
		Accounts: map[string]account.Account{
			"exchange": {
				Balances: map[string]float64{"DAI": 1000},
			},
		},
	})

	// Check that returned value is correct
	suite.Require().Equal(appSetID, id)
	suite.Require().NoError(err)
}
