package exchanges

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/exchanges"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
)

func TestGetCachedSuite(t *testing.T) {
	suite.Run(t, new(GetCachedSuite))
}

type GetCachedSuite struct {
	suite.Suite
	operator Operator
	db       *db.MockAdapter
	exchange *exchanges.MockAdapter
}

func (suite *GetCachedSuite) SetupTest() {
	suite.db = db.NewMockAdapter(gomock.NewController(suite.T()))

	suite.exchange = exchanges.NewMockAdapter(gomock.NewController(suite.T()))
	exchanges := map[string]exchanges.Adapter{"exchange": suite.exchange}

	suite.operator = New(suite.db, exchanges)
}

func (suite *GetCachedSuite) setMocksForReadExchangesThatIsNotCached() context.Context {
	ctx := context.Background()

	// Set call to DB that should not return anything
	suite.db.EXPECT().ReadExchanges(ctx, "exchange").Return(
		[]exchange.Exchange{},
		nil,
	)

	// Set call to exchange that should return info
	suite.exchange.EXPECT().Infos(ctx).Return(
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
	exchanges, err := suite.operator.GetCached(ctx, DurationOpt(time.Hour), "exchange")

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
