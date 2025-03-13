package activities

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	_ "github.com/lib/pq" // PostGres driver
)

// SQL is a struct that contains the SQL connection.
type SQL struct {
	DB *sqlx.DB
}

// NewSQL creates a new SQL connection.
func NewSQL(ctx context.Context, c config.SQL) (SQL, error) {
	if err := c.Validate(); err != nil {
		return SQL{}, err
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", c.DSN)
	if err != nil {
		return SQL{}, err
	}

	return SQL{
		DB: db,
	}, nil
}
