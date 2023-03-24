package ticks

import "github.com/digital-feather/cryptellation/pkg/tick"

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.events.Subscribe(pairSymbol)
}
