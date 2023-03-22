package app

import (
	"time"

	"github.com/digital-feather/cryptellation/pkg/types/period"
)

type GetCachedPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}
