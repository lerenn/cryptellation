// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

type Port interface {
	GetSMA(ctx context.Context, payload ReadSMAPayload) (*timeserie.TimeSerie[float64], error)
	UpsertSMA(ctx context.Context, payload WriteSMAPayload) error
}

type ReadSMAPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	PeriodNumber uint
	PriceType    candlestick.PriceType
	Start        time.Time
	End          time.Time
}

type WriteSMAPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	PeriodNumber uint
	PriceType    candlestick.PriceType
	TimeSerie    *timeserie.TimeSerie[float64]
}
