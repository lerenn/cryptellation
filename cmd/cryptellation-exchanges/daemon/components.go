package daemon

import "github.com/lerenn/cryptellation/internal/components/exchanges"

type components struct {
	exchanges exchanges.Interface
}

func newComponents(adapters adapters) components {
	return components{
		exchanges: exchanges.New(adapters.db, adapters.exchanges),
	}
}

func (c components) Close() {
}
