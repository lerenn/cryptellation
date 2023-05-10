package event

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/candlestick"
	"github.com/lerenn/cryptellation/pkg/tick"
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
		Time:       t,
		PairSymbol: pairSymbol,
		Price:      cs.PriceByType(currentPriceType),
		Exchange:   exchange,
	}), nil
}
