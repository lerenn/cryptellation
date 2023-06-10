// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import "github.com/lerenn/cryptellation/pkg/models/tick"

type Port interface {
	ListenSymbol(exchange, symbol string) (chan tick.Tick, chan struct{}, error)
}
