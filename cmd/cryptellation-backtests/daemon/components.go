package daemon

import "github.com/lerenn/cryptellation/internal/components/backtests"

type components struct {
	backtests backtests.Interface
}

func newComponents(adapters adapters) components {
	return components{
		backtests: backtests.New(adapters.db, adapters.events, adapters.candlesticks),
	}
}

func (c components) Close() {
}
