package tests

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/services/exchanges/internal/application/ports/db"
	"github.com/digital-feather/cryptellation/services/exchanges/pkg/models/exchange"
	"github.com/stretchr/testify/suite"
)

type ExchangesSuite struct {
	suite.Suite
	DB db.Adapter
}

func (suite *ExchangesSuite) TestCreateRead() {
	as := suite.Require()

	// Given a exchange
	p := exchange.Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-DEF", "IJK-LMN"},
		PeriodsSymbols: []string{"M1", "M3"},
		Fees:           0.1,
		LastSyncTime:   time.Now().UTC(),
	}

	// When we create it and read it
	as.NoError(suite.DB.CreateExchanges(context.Background(), p))
	rp, err := suite.DB.ReadExchanges(context.Background(), p.Name)
	as.NoError(err)

	// Then it's the same
	as.Len(rp, 1)
	as.Equal(p.Name, rp[0].Name)
	as.Contains(rp[0].PairsSymbols, "ABC-DEF")
	as.Contains(rp[0].PairsSymbols, "IJK-LMN")
	as.Contains(rp[0].PeriodsSymbols, "M1")
	as.Contains(rp[0].PeriodsSymbols, "M3")
	as.Equal(p.Fees, rp[0].Fees)
	as.WithinDuration(time.Now().UTC(), rp[0].LastSyncTime, time.Second)
}

func (suite *ExchangesSuite) TestCreateReadInexistant() {
	as := suite.Require()

	// When we read an inexistant exchange
	exchanges, err := suite.DB.ReadExchanges(context.Background(), "inexistant")

	// Then there is no error but no exchange
	as.NoError(err)
	as.Len(exchanges, 0)
}

func (suite *ExchangesSuite) TestReadAll() {
	as := suite.Require()

	// Given 3 created exchanges
	p1 := exchange.Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-DEF", "IJK-LMN"},
		PeriodsSymbols: []string{"M1", "M3"},
		Fees:           0.1,
		LastSyncTime:   time.Now().UTC(),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p1))
	p2 := exchange.Exchange{
		Name:           "exchange2",
		PairsSymbols:   []string{"ABC-DEF", "XYZ-LMN"},
		PeriodsSymbols: []string{"M1", "M5"},
		Fees:           0.2,
		LastSyncTime:   time.Now().UTC(),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p2))
	p3 := exchange.Exchange{
		Name:           "exchange3",
		PairsSymbols:   []string{"ABC-XYZ", "IJK-LMN"},
		PeriodsSymbols: []string{"M1", "M10"},
		Fees:           0.3,
		LastSyncTime:   time.Now().UTC(),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p3))

	// When we read all of them
	ps, err := suite.DB.ReadExchanges(context.Background())
	as.NoError(err)

	// Then we have all of them
	as.Len(ps, 3)
	for _, p := range ps {
		if p.Name != p1.Name && p.Name != p2.Name && p.Name != p3.Name {
			as.Fail("This exchange should not exists", p)
		}
	}
}

func (suite *ExchangesSuite) TestUpdate() {
	as := suite.Require()

	// Given a created exchange
	p1 := exchange.Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-DEF", "IJK-LMN"},
		PeriodsSymbols: []string{"M1", "M3"},
		Fees:           0.1,
		LastSyncTime:   time.Now().UTC().Add(-time.Hour),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p1))

	// When we make some modification to it
	p2 := exchange.Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-XYZ", "IJK-XYZ"},
		PeriodsSymbols: []string{"M5", "D1"},
		Fees:           0.2,
		LastSyncTime:   time.Now().UTC(),
	}

	// And we update it
	as.NoError(suite.DB.UpdateExchanges(context.Background(), p2))

	// Then the exchange is updated
	rp, err := suite.DB.ReadExchanges(context.Background(), p2.Name)
	as.NoError(err)
	as.Len(rp, 1)
	as.Equal(p2.Name, rp[0].Name)
	as.Contains(rp[0].PairsSymbols, "ABC-XYZ")
	as.Contains(rp[0].PairsSymbols, "IJK-XYZ")
	as.Contains(rp[0].PeriodsSymbols, "M5")
	as.Contains(rp[0].PeriodsSymbols, "D1")
	as.Equal(p2.Fees, rp[0].Fees)
	as.WithinDuration(time.Now().UTC(), rp[0].LastSyncTime, time.Second)
}

func (suite *ExchangesSuite) TestDelete() {
	as := suite.Require()

	// Given twp created exchange
	p1 := exchange.Exchange{
		Name:           "exchange",
		PairsSymbols:   []string{"ABC-XYZ", "IJK-XYZ"},
		PeriodsSymbols: []string{"M5", "D1"},
		Fees:           0.2,
		LastSyncTime:   time.Now().UTC(),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p1))
	p2 := exchange.Exchange{
		Name:           "exchange2",
		PairsSymbols:   []string{"ABC-XYZ", "IJK-ABC"},
		PeriodsSymbols: []string{"M5", "M1"},
		Fees:           0.3,
		LastSyncTime:   time.Now().UTC(),
	}
	as.NoError(suite.DB.CreateExchanges(context.Background(), p2))

	// When we delete it
	as.NoError(suite.DB.DeleteExchanges(context.Background(), p1.Name))

	// Then we can't read it anymore
	exchanges, err := suite.DB.ReadExchanges(context.Background(), p1.Name, p2.Name)
	as.NoError(err)
	as.Len(exchanges, 1)
	as.Equal(p2.Name, exchanges[0].Name)
}
