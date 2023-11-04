package daemon

import "github.com/lerenn/cryptellation/internal/components/indicators"

type components struct {
	indicators indicators.Interface
}

func newComponents(adapters adapters) components {
	return components{
		indicators: indicators.New(adapters.db, adapters.candlesticks),
	}
}

func (c components) Close() {
}
