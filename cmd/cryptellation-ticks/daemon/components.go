package daemon

import "github.com/lerenn/cryptellation/internal/components/ticks"

type components struct {
	ticks ticks.Interface
}

func newComponents(adapters adapters) components {
	return components{
		ticks: ticks.New(adapters.events, adapters.db, adapters.exchanges),
	}
}

func (c components) Close() {
}
