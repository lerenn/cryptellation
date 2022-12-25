package ticks

import "github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.pubsub.Subscribe(pairSymbol)
}
