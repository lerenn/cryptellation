// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=indicators.go -destination=mock/indicators.gen.go -package mock

package client

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

type Indicators interface {
	SMA(ctx context.Context, payload SMAPayload) (*timeserie.TimeSerie[float64], error)
	Close()
}

type SMAPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        time.Time
	End          time.Time
	PeriodNumber uint
	PriceType    candlestick.PriceType
}
