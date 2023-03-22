package app

import "github.com/digital-feather/cryptellation/pkg/types/tick"

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.events.Subscribe(pairSymbol)
}
