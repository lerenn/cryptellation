package sql

import (
	"context"
	"testing"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/indicators/deployments"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/ports/db"
	"github.com/stretchr/testify/suite"
)

func TestIndicatorsSuite(t *testing.T) {
	suite.Run(t, new(IndicatorsSuite))
}

type IndicatorsSuite struct {
	db.IndicatorsSuite
}

func (suite *IndicatorsSuite) SetupTest() {
	db, err := New(config.LoadSQL(&config.SQL{
		Database: "indicators",
		Port:     deployments.DockerComposeSQLDBPort,
	}))
	suite.Require().NoError(err)
	suite.Require().NoError(db.Reset(context.TODO()))

	suite.DB = db
}
