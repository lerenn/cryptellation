// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=adapter.go -destination=mock.gen.go -package exchanges

package exchanges

import "github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"

type Adapter interface {
	ListenSymbol(symbol string) (chan tick.Tick, chan struct{}, error)
	Name() string
}
