package daemon

import (
	"cryptellation/svc/forwardtests/internal/app"
	"cryptellation/svc/forwardtests/internal/app/domain"
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
