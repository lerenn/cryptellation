package daemon

import "github.com/lerenn/cryptellation/internal/core/candlesticks"

type components struct {
	candlesticks candlesticks.Interface
}

func newComponents(adapters adapters) components {
	return components{
		candlesticks: candlesticks.New(adapters.db, adapters.exchanges),
	}
}

func (c components) Close() {
}
