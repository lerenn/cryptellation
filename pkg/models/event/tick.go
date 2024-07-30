package event

import (
	"time"

	"cryptellation/svc/candlesticks/pkg/candlestick"

	"cryptellation/svc/ticks/pkg/tick"
)

func NewTickEvent(t time.Time, content tick.Tick) Event {
	return Event{
		Type:    TypeIsTick,
		Time:    t,
		Content: content,
	}
}

func TickEventFromCandlestick(
	exchange, pair string,
	currentPriceType candlestick.PriceType,
	t time.Time,
	cs candlestick.Candlestick,
) (Event, error) {
	return NewTickEvent(t, tick.Tick{
		Time:     t,
		Pair:     pair,
		Price:    cs.PriceByType(currentPriceType),
		Exchange: exchange,
	}), nil
}
