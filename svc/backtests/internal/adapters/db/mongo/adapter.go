package backtests

import (
	mongoutil "github.com/lerenn/cryptellation/pkg/adapters/db/mongo"
	"github.com/lerenn/cryptellation/pkg/config"
	port "github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
)

var _ port.Port = (*Adapter)(nil)

type Adapter struct {
	client mongoutil.Client
}

func New(c config.Mongo) (*Adapter, error) {
	// Create embedded database access
	db, err := mongoutil.NewClient(c)

	// Return database access
	return &Adapter{
		client: db,
	}, err
}
