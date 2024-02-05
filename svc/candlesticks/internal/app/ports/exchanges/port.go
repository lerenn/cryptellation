// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type Port interface {
	GetCandlesticks(ctx context.Context, payload GetCandlesticksPayload) (*candlestick.List, error)
}

type GetCandlesticksPayload struct {
	Exchange string
	Pair     string
	Period   period.Symbol
	Start    time.Time
	End      time.Time
	Limit    int
}
