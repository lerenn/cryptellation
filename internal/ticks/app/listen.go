package app

import "github.com/digital-feather/cryptellation/pkg/tick"

func (t Ticks) Listen(exchange, pairSymbol string) (<-chan tick.Tick, error) {
	return t.pubsub.Subscribe(pairSymbol)
}
