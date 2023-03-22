package exchanges

import (
	"time"

	"github.com/digital-feather/cryptellation/pkg/types/period"
)

type GetCandlesticksPayload struct {
	Exchange   string
	PairSymbol string
	Period     period.Symbol
	Start      time.Time
	End        time.Time
	Limit      int
}
