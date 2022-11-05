package exchanges

import (
	"time"

	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/period"
)

type GetCandlesticksPayload struct {
	PairSymbol string
	Period     period.Symbol
	Start      time.Time
	End        time.Time
	Limit      int
}
