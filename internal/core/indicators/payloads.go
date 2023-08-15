package indicators

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/pkg/models/period"
)

type GetCachedSMAPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        time.Time
	End          time.Time
	PeriodNumber uint
	PriceType    candlestick.PriceType
}
