package activities

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	_ "github.com/lib/pq" // PostGres driver
)

// PostGres is a struct that contains the PostGres connection.
type PostGres struct {
	DB *sqlx.DB
}

// NewPostGres creates a new PostGres connection.
func NewPostGres(ctx context.Context, c config.PostGres) (PostGres, error) {
	if err := c.Validate(); err != nil {
		return PostGres{}, err
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", c.DSN)
	if err != nil {
		return PostGres{}, err
	}

	return PostGres{
		DB: db,
	}, nil
}
