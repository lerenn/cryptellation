package event

import (
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/tick"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

func NewTickEvent(t time.Time, content tick.Tick) Event {
	return Event{
		Type:    TypeIsTick,
		Time:    t,
		Content: content,
	}
}

func TickEventFromCandlestick(
	exchange, pairSymbol string,
	currentPriceType candlestick.PriceType,
	t time.Time,
	cs candlestick.Candlestick,
) (Event, error) {
	return NewTickEvent(t, tick.Tick{
		PairSymbol: pairSymbol,
		Price:      cs.PriceByType(currentPriceType),
		Exchange:   exchange,
	}), nil
}
