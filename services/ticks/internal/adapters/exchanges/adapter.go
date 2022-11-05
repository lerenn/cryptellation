package exchanges

import "github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"

type Adapter interface {
	ListenSymbol(symbol string) (chan tick.Tick, chan struct{}, error)
}
