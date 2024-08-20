package daemon

import (
	"github.com/lerenn/cryptellation/svc/ticks/internal/app"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/domain"
)

type components struct {
	ticks app.Ticks
}

func newComponents(adapters adapters) components {
	return components{
		ticks: domain.New(adapters.exchanges, adapters.events),
	}
}

func (c components) Close() {
}
