package daemon

import (
	"github.com/lerenn/cryptellation/svc/backtests/internal/app"
	"github.com/lerenn/cryptellation/svc/backtests/internal/app/domain"
)

type components struct {
	backtests app.Backtests
}

func newComponents(adapters adapters) components {
	return components{
		backtests: domain.New(adapters.db, adapters.events, adapters.candlesticks),
	}
}

func (c components) Close() {
}
