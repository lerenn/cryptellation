package candlesticks

import (
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
)

type GetCachedPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}
