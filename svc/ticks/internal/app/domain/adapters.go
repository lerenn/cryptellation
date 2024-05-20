package domain

import (
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/events"
	"github.com/lerenn/cryptellation/svc/ticks/internal/app/ports/exchanges"
)

type adapters struct {
	Exchanges exchanges.Port
	Events    events.Port
}
