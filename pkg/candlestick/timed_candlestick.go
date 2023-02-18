package candlestick

import "time"

type TimedCandlestick struct {
	Time time.Time
	Candlestick
}
