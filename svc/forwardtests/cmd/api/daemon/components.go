package daemon

import (
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/domain"
)

type components struct {
	forwardtests app.ForwardTests
}

func newComponents(adapters adapters) components {
	return components{
		forwardtests: domain.New(adapters.db, adapters.candlesticks),
	}
}

func (c components) Close() {
}
