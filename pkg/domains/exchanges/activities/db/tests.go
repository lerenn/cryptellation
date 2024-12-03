package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"github.com/stretchr/testify/suite"
)

// ExchangesSuite is the test suite for the exchanges database.
type ExchangesSuite struct {
	suite.Suite
	DB DB
}

// TestCreateRead will test the creation and reading of an exchange.
func (suite *ExchangesSuite) TestCreateRead() {
	// Given a exchange
	p := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC(),
	}

	// When we create it and read it
	_, err := suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p},
		})
	suite.Require().NoError(err)
	r, err := suite.DB.ReadExchangesActivity(context.Background(), ReadExchangesActivityParams{
		Names: []string{p.Name},
	})
	suite.Require().NoError(err)

	// Then it's the same
	suite.Require().Len(r.Exchanges, 1)
	suite.Require().Equal(p.Name, r.Exchanges[0].Name)
	suite.Require().Contains(r.Exchanges[0].Pairs, "ABC-DEF")
	suite.Require().Contains(r.Exchanges[0].Pairs, "IJK-LMN")
	suite.Require().Contains(r.Exchanges[0].Periods, "M1")
	suite.Require().Contains(r.Exchanges[0].Periods, "M3")
	suite.Require().Equal(p.Fees, r.Exchanges[0].Fees)
	suite.Require().WithinDuration(time.Now().UTC(), r.Exchanges[0].LastSyncTime, time.Second)
}

// TestCreateReadInexistant will test the reading of an inexistant exchange.
func (suite *ExchangesSuite) TestCreateReadInexistant() {
	// When we read an inexistant exchange
	r, err := suite.DB.ReadExchangesActivity(context.Background(), ReadExchangesActivityParams{
		Names: []string{"inexistant"},
	})

	// Then there is no error but no exchange
	suite.Require().NoError(err)
	suite.Require().Len(r.Exchanges, 0)
}

// TestReadAll will test the reading of all exchanges.
func (suite *ExchangesSuite) TestReadAll() {
	// Given 3 created exchanges
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC(),
	}
	_, err := suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p1},
		})
	suite.Require().NoError(err)
	p2 := exchange.Exchange{
		Name:         "exchange2",
		Pairs:        []string{"ABC-DEF", "XYZ-LMN"},
		Periods:      []string{"M1", "M5"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	_, err = suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p2},
		})
	suite.Require().NoError(err)
	p3 := exchange.Exchange{
		Name:         "exchange3",
		Pairs:        []string{"ABC-XYZ", "IJK-LMN"},
		Periods:      []string{"M1", "M10"},
		Fees:         0.3,
		LastSyncTime: time.Now().UTC(),
	}
	_, err = suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p3},
		})
	suite.Require().NoError(err)

	// When we read all of them
	r, err := suite.DB.ReadExchangesActivity(context.Background(), ReadExchangesActivityParams{})
	suite.Require().NoError(err)

	// Then we have all of them
	suite.Require().Len(r.Exchanges, 3)
	for _, p := range r.Exchanges {
		if p.Name != p1.Name && p.Name != p2.Name && p.Name != p3.Name {
			suite.Require().Fail("This exchange should not exists", p)
		}
	}
}

// TestUpdate will test the update of an exchange.
func (suite *ExchangesSuite) TestUpdate() {
	// Given a created exchange
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-DEF", "IJK-LMN"},
		Periods:      []string{"M1", "M3"},
		Fees:         0.1,
		LastSyncTime: time.Now().UTC().Add(-time.Hour),
	}
	_, err := suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p1},
		})
	suite.Require().NoError(err)

	// When we make some modification to it
	p2 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}

	// And we update it
	_, err = suite.DB.UpdateExchangesActivity(context.Background(), UpdateExchangesActivityParams{
		Exchanges: []exchange.Exchange{p2},
	})
	suite.Require().NoError(err)

	// Then the exchange is updated
	r, err := suite.DB.ReadExchangesActivity(context.Background(), ReadExchangesActivityParams{
		Names: []string{p2.Name},
	})
	suite.Require().NoError(err)
	suite.Require().Len(r.Exchanges, 1)
	suite.Require().Equal(p2.Name, r.Exchanges[0].Name)
	suite.Require().Contains(r.Exchanges[0].Pairs, "ABC-XYZ")
	suite.Require().Contains(r.Exchanges[0].Pairs, "IJK-XYZ")
	suite.Require().Contains(r.Exchanges[0].Periods, "M5")
	suite.Require().Contains(r.Exchanges[0].Periods, "D1")
	suite.Require().Equal(p2.Fees, r.Exchanges[0].Fees)
	suite.Require().WithinDuration(time.Now().UTC(), r.Exchanges[0].LastSyncTime, time.Second)
}

// TestDelete will test the deletion of an exchange.
func (suite *ExchangesSuite) TestDelete() {
	// Given twp created exchange
	p1 := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	_, err := suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p1},
		})
	suite.Require().NoError(err)
	p2 := exchange.Exchange{
		Name:         "exchange2",
		Pairs:        []string{"ABC-XYZ", "IJK-ABC"},
		Periods:      []string{"M5", "M1"},
		Fees:         0.3,
		LastSyncTime: time.Now().UTC(),
	}
	_, err = suite.DB.CreateExchangesActivity(context.Background(),
		CreateExchangesActivityParams{
			Exchanges: []exchange.Exchange{p2},
		})
	suite.Require().NoError(err)

	// When we delete it
	_, err = suite.DB.DeleteExchangesActivity(context.Background(), DeleteExchangesActivityParams{
		Names: []string{p1.Name},
	})
	suite.Require().NoError(err)

	// Then we can't read it anymore
	r, err := suite.DB.ReadExchangesActivity(context.Background(), ReadExchangesActivityParams{
		Names: []string{p1.Name, p2.Name},
	})
	suite.Require().NoError(err)
	suite.Require().Len(r.Exchanges, 1)
	suite.Require().Equal(p2.Name, r.Exchanges[0].Name)
}
