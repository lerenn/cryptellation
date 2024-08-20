package daemon

import (
	"github.com/lerenn/cryptellation/candlesticks/internal/app"
	"github.com/lerenn/cryptellation/candlesticks/internal/app/domain"
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
