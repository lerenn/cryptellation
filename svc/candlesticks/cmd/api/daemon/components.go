package daemon

import (
	"cryptellation/svc/candlesticks/internal/app"
	"cryptellation/svc/candlesticks/internal/app/domain"
)

type components struct {
	candlesticks app.Candlesticks
}

func newComponents(adapters adapters) components {
	return components{
		candlesticks: domain.New(adapters.db, adapters.exchanges),
	}
}

func (c components) Close() {
}
