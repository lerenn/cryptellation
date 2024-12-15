package entities

import (
	"testing"
	"time"

	binance "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/v1/pkg/activities"
	"github.com/lerenn/cryptellation/v1/pkg/models/candlestick"
	"github.com/lerenn/cryptellation/v1/pkg/models/period"
)

var testCasesKLineToCandlestick = []struct {
	KLine       binance.Kline
	Candlestick candlestick.Candlestick
}{
	{
		KLine: binance.Kline{
			OpenTime: 0,
			Open:     "1.0", High: "2.0", Low: "0.5", Close: "1.5", Volume: "1000"},
		Candlestick: candlestick.Candlestick{
			Time: time.Unix(0, 0),
			Open: 1, High: 2, Low: 0.5, Close: 1.5, Volume: 1000, Uncomplete: false},
	},
	{
		KLine: binance.Kline{
			OpenTime: 60 * 1000,
			Open:     "2.0", High: "4.0", Low: "1", Close: "3", Volume: "1500"},
		Candlestick: candlestick.Candlestick{
			Time: time.Unix(60, 0),
			Open: 2, High: 4, Low: 1, Close: 3, Volume: 1500, Uncomplete: false},
	},
}

func TestKLineToCandlestick(t *testing.T) {
	for i, test := range testCasesKLineToCandlestick {
		cs, err := KLineToCandlestick(test.KLine, period.M1, time.Unix(120, 0))
		switch {
		case err != nil:
			t.Error("There should be no error on Candlestick", i, ":", err)
		case test.Candlestick != cs:
			t.Error("Candlestick", i, "is not transformed correctly:", test.Candlestick, cs)
		case !test.Candlestick.Time.Equal(cs.Time):
			t.Error("times should be equal")
		}
	}
}

func TestKLineToCandlestick_IncorrectOpen(t *testing.T) {
	c := binance.Kline{OpenTime: 0, Open: "error", High: "2.0", Low: "0.5", Close: "1.5"}
	if _, err := KLineToCandlestick(c, period.M1, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on open")
	}
}

func TestKLineToCandlestick_IncorrectHigh(t *testing.T) {
	c := binance.Kline{OpenTime: 0, Open: "1.0", High: "error", Low: "0.5", Close: "1.5"}
	if _, err := KLineToCandlestick(c, period.M1, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on high")
	}
}

func TestKLineToCandlestick_IncorrectLow(t *testing.T) {
	c := binance.Kline{OpenTime: 0, Open: "1.0", High: "2.0", Low: "error", Close: "1.5"}
	if _, err := KLineToCandlestick(c, period.M1, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on low")
	}
}

func TestKLineToCandlestick_IncorrectClose(t *testing.T) {
	c := binance.Kline{OpenTime: 0, Open: "1.0", High: "2.0", Low: "0.5", Close: "error"}
	if _, err := KLineToCandlestick(c, period.M1, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on close")
	}
}

func TestKLinesToCandlesticks(t *testing.T) {
	p := "BTC-USDC"

	// Only get klines
	kl := make([]*binance.Kline, len(testCasesKLineToCandlestick))
	for i := range testCasesKLineToCandlestick {
		kl[i] = &testCasesKLineToCandlestick[i].KLine
	}

	// Test function
	cs, err := KLinesToCandlesticks(p, period.M1, kl, time.Unix(120, 0))
	if err != nil {
		t.Error("There should be no error:", err)
	}

	switch {
	case cs.Metadata.Exchange != activities.BinanceInfos.Name:
		t.Fatal("Exchange should be binance, but is", cs.Metadata.Exchange)
	case cs.Metadata.Pair != p:
		t.Fatal("Pair should be", p, "but is", cs.Metadata.Pair)
	case cs.Metadata.Period != period.M1:
		t.Fatal("Period should be", period.M1, "but is", cs.Metadata.Period)
	}

	for i, test := range testCasesKLineToCandlestick {
		rc, _ := cs.Data.Get(test.Candlestick.Time)
		if test.Candlestick != rc {
			t.Error("Candlestick", i, "is not transformed correctly:", test.Candlestick, rc)
		}
	}
}

func TestKLinesToCandlesticks_IncorrectOpen(t *testing.T) {
	c := []*binance.Kline{{OpenTime: 0, Open: "error", High: "2.0", Low: "0.5", Close: "1.5", Volume: "1000"}}
	if _, err := KLinesToCandlesticks("BTC-USDC", period.M1, c, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on open")
	}
}

func TestKLinesToCandlesticks_IncorrectHigh(t *testing.T) {
	c := []*binance.Kline{{OpenTime: 0, Open: "1.0", High: "error", Low: "0.5", Close: "1.5", Volume: "1000"}}
	if _, err := KLinesToCandlesticks("BTC-USDC", period.M1, c, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on high")
	}
}

func TestKLinesToCandlesticks_IncorrectLow(t *testing.T) {
	c := []*binance.Kline{{OpenTime: 0, Open: "1.0", High: "2.0", Low: "error", Close: "1.5", Volume: "1000"}}
	if _, err := KLinesToCandlesticks("BTC-USDC", period.M1, c, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on low")
	}
}

func TestKLinesToCandlesticks_IncorrectClose(t *testing.T) {
	c := []*binance.Kline{{OpenTime: 0, Open: "1.0", High: "2.0", Low: "0.5", Close: "error", Volume: "1000"}}
	if _, err := KLinesToCandlesticks("BTC-USDC", period.M1, c, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on close")
	}
}

func TestKLinesToCandlesticks_IncorrectVolume(t *testing.T) {
	c := []*binance.Kline{{OpenTime: 0, Open: "1.0", High: "2.0", Low: "0.5", Close: "1.5", Volume: "error"}}
	if _, err := KLinesToCandlesticks("BTC-USDC", period.M1, c, time.Unix(120, 0)); err == nil {
		t.Error("There should be an error on close")
	}
}

var KLineTimeTotimeTests = []struct {
	Binancetime int64
	Time        time.Time
}{
	{Binancetime: 1257894000000, Time: time.Unix(1257894000, 0)},
}

func TestKLineTimeToTime(t *testing.T) {
	for i, c := range KLineTimeTotimeTests {
		r := KLineTimeToTime(c.Binancetime)
		if !r.Equal(c.Time) {
			t.Error("Times don't match on test", i, ":", c.Time, r)
		}
	}
}

func TestTimeToKLineTime(t *testing.T) {
	for i, c := range KLineTimeTotimeTests {
		r := TimeToKLineTime(c.Time)
		if r != c.Binancetime {
			t.Error("Times don't match on test", i, ":", c.Binancetime, r)
		}
	}
}

func TestKLineToCandlestick_Uncomplete(t *testing.T) {
	c := binance.Kline{OpenTime: 120 * 1000, Open: "1.0", High: "2.0", Low: "0.5", Close: "1.5", Volume: "1000"}
	cs, err := KLineToCandlestick(c, period.M1, time.Unix(130, 0))
	if err != nil {
		t.Error("There should be no error:", err.Error())
	}

	if !cs.Uncomplete {
		t.Error("Candlestick should be uncomplete")
	}
}
