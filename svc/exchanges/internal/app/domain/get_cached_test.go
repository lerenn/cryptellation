package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/lerenn/cryptellation/svc/exchanges/internal/app"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/app/ports/exchanges"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
)

func TestGetCachedSuite(t *testing.T) {
	suite.Run(t, new(GetCachedSuite))
}

type GetCachedSuite struct {
	suite.Suite
	app      app.Exchanges
	db       *db.MockPort
	exchange *exchanges.MockPort
}

func (suite *GetCachedSuite) SetupTest() {
	suite.db = db.NewMockPort(gomock.NewController(suite.T()))
	suite.exchange = exchanges.NewMockPort(gomock.NewController(suite.T()))

	suite.app = New(suite.db, suite.exchange)
}

func (suite *GetCachedSuite) setMocksForReadExchangesThatIsNotCached() context.Context {
	ctx := context.Background()

	// Set call to DB that should not return anything
	suite.db.EXPECT().ReadExchanges(ctx, "exchange").Return(
		[]exchange.Exchange{},
		nil,
	)

	// Set call to exchange that should return info
	suite.exchange.EXPECT().Infos(ctx, "exchange").Return(
		exchange.Exchange{Name: "exchange"},
		nil,
	)

	// Set call to database that should create the exchange
	suite.db.EXPECT().CreateExchanges(ctx, exchange.Exchange{Name: "exchange"}).Return(nil)

	return ctx
}

func (suite *GetCachedSuite) TestReadExchangesThatIsNotCached() {
	ctx := suite.setMocksForReadExchangesThatIsNotCached()

	// When requesting an exchange for the first time
	exchanges, err := suite.app.GetCached(ctx, "exchange")

	// Then the request is successful
	suite.Require().NoError(err)

	// And the exchange is correct
	suite.Require().Len(exchanges, 1)
	suite.Require().Equal("exchange", exchanges[0].Name)
}

func (suite *GetCachedSuite) TestReadExchangesThatIsCached() {
	// TODO
}

func (suite *GetCachedSuite) TestReadExchangesThatIsOutdatedCached() {
	// TODO
}

func DurationOpt(t time.Duration) *time.Duration {
	return &t
}
