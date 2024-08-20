// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package db

package db

import (
	"context"
	"time"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"

	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/candlesticks/pkg/period"
)

type Port interface {
	GetSMA(ctx context.Context, payload ReadSMAPayload) (*timeserie.TimeSerie[float64], error)
	UpsertSMA(ctx context.Context, payload WriteSMAPayload) error
}

type ReadSMAPayload struct {
	Exchange     string
	Pair         string
	Period       period.Symbol
	PeriodNumber int
	PriceType    candlestick.PriceType
	Start        time.Time
	End          time.Time
}

type WriteSMAPayload struct {
	Exchange     string
	Pair         string
	Period       period.Symbol
	PeriodNumber int
	PriceType    candlestick.PriceType
	TimeSerie    *timeserie.TimeSerie[float64]
}
