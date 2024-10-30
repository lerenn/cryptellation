package event

import (
	"time"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func NewPriceEvent(t time.Time, content tick.Tick) Event {
	return Event{
		Type:    TypeIsPrice,
		Time:    t,
		Content: content,
	}
}

func PriceEventFromCandlestick(
	exchange, pair string,
	currentPriceType candlestick.Price,
	t time.Time,
	cs candlestick.Candlestick,
) (Event, error) {
	return NewPriceEvent(t, tick.Tick{
		Time:     t,
		Pair:     pair,
		Price:    cs.Price(currentPriceType),
		Exchange: exchange,
	}), nil
}
