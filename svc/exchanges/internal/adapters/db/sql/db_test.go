package exchanges

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/exchanges/internal/adapters/db/sql/entities"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
	"github.com/stretchr/testify/suite"
)

func TestSqlDatabaseSuite(t *testing.T) {
	suite.Run(t, new(SqlDatabaseSuite))
}

type SqlDatabaseSuite struct {
	suite.Suite
	adapter *Adapter
}

func (suite *SqlDatabaseSuite) SetupTest() {
	defer tmpEnvVar("SQLDB_DATABASE", "exchanges")()

	db, err := New(config.LoadSQL())
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.adapter = db
}

func (suite *SqlDatabaseSuite) TestNewWithURIError() {
	defer tmpEnvVar("SQLDB_HOST", "")()

	var err error
	_, err = New(config.LoadSQL())
	suite.Error(err)
}

func (suite *SqlDatabaseSuite) TestReset() {
	as := suite.Require()

	// Given a created exchange
	p := exchange.Exchange{
		Name:         "exchange",
		Pairs:        []string{"ABC-XYZ", "IJK-XYZ"},
		Periods:      []string{"M5", "D1"},
		Fees:         0.2,
		LastSyncTime: time.Now().UTC(),
	}
	as.NoError(suite.adapter.CreateExchanges(context.Background(), p))

	// When we reset the DB
	defer tmpEnvVar("SQLDB_DATABASE", "exchanges")()
	as.NoError(suite.adapter.Reset(context.TODO()))

	// Then there is no exchange left
	exchanges := []entities.Exchange{}
	as.NoError(suite.adapter.db.Client.WithContext(context.Background()).Find(&exchanges).Error)
	as.Len(exchanges, 0)

	// And there is no pair left
	pairs := []entities.Pair{}
	as.NoError(suite.adapter.db.Client.WithContext(context.Background()).Find(&pairs).Error)
	as.Len(pairs, 0)

	// And there is no period left
	periods := []entities.Period{}
	as.NoError(suite.adapter.db.Client.WithContext(context.Background()).Find(&periods).Error)
	as.Len(periods, 0)
}

func tmpEnvVar(key, value string) (reset func()) {
	originalValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, originalValue)
	}
}
