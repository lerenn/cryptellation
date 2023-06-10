package ticks

import "github.com/lerenn/cryptellation/pkg/models/tick"

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.events.Subscribe(pairSymbol)
}
