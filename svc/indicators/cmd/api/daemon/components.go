package daemon

import (
	"github.com/lerenn/cryptellation/svc/indicators/internal/app"
	"github.com/lerenn/cryptellation/svc/indicators/internal/app/domain"
)

type components struct {
	indicators app.Indicators
}

func newComponents(adapters adapters) components {
	return components{
		indicators: domain.New(adapters.db, adapters.candlesticks),
	}
}

func (c components) Close() {
}
