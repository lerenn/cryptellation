package domain

import (
	"cryptellation/svc/ticks/internal/app/ports/events"
	"cryptellation/svc/ticks/internal/app/ports/exchanges"
)

type adapters struct {
	Exchanges exchanges.Port
	Events    events.Port
}
