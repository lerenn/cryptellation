package candlesticks

import (
	"math"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
)

const (
	unicodeVoid           = " "
	unicodeBody           = "┃"
	unicodeHalfBodyTop    = "╻"
	unicodeHalfBodyBottom = "╹"
	unicodeWick           = "│"
	unicodeTop            = "╽"
	unicodeBottom         = "╿"
	unicodeUpperWick      = "╷"
	unicodeLowerWick      = "╵"

	unicodeFlat                      = "-"
	unicodeMinimalBodyFullWick       = "┿"
	unicodeMinimalBodyTopHalfWick    = "┷"
	unicodeMinimalBodyBottomHalfWick = "┯"
)

type column struct {
	symbols []string
	isUp    bool
}

func newColumn(c candlestick.Candlestick, min, max float64, height int) column {
	bodyTop := c.Open
	bodyBottom := c.Close
	if c.Open <= c.Close {
		bodyTop = c.Close
		bodyBottom = c.Open
	}
	wickTop := c.High
	wickBottom := c.Low

	nBodyTop := normalizePoint(bodyTop, min, max, height)
	nBodyBottom := normalizePoint(bodyBottom, min, max, height)
	nWickTop := normalizePoint(wickTop, min, max, height)
	nWickBottom := normalizePoint(wickBottom, min, max, height)

	rBodyTop := math.Round(nBodyTop)
	rBodyBottom := math.Round(nBodyBottom)
	rWickTop := math.Round(nWickTop)
	rWickBottom := math.Round(nWickBottom)

	qrBodyTop := math.Round(nBodyTop*2) / 2
	qrBodyBottom := math.Round(nBodyBottom*2) / 2
	qrWickTop := math.Round(nWickTop*2) / 2
	qrWickBottom := math.Round(nWickBottom*2) / 2

	symbols := make([]string, height)
	// i == 0 == top of the screen (i.e. reversed y axis)
	for i := float64(0); i < float64(height); i++ {
		var symbol string

		switch {
		case rWickTop < i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			symbol = unicodeVoid
		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom > i:
			symbol = unicodeVoid

		// BOTTOM

		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom == i:
			if qrWickBottom-i == 0 {
				symbol = unicodeWick
			} else {
				symbol = unicodeLowerWick
			}

		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom < i:
			symbol = unicodeWick

		case rWickTop > i && rBodyTop > i && rBodyBottom == i && rWickBottom == i:
			if qrBodyBottom == 0 {
				symbol = unicodeBody
			} else if qrWickBottom == 0 {
				symbol = unicodeBottom
			} else {
				symbol = unicodeHalfBodyBottom
			}

		// MIDDLE

		case rWickTop > i && rBodyTop > i && rBodyBottom == i && rWickBottom < i:
			if qrBodyBottom == 0 {
				symbol = unicodeBody
			} else {
				symbol = unicodeBottom
			}

		case rWickTop > i && rBodyTop > i && rBodyBottom < i && rWickBottom < i:
			symbol = unicodeBody

		case rWickTop > i && rBodyTop == i && rBodyBottom == i && rWickBottom < i:
			if qrBodyTop == 0 {
				symbol = unicodeWick
			} else if qrBodyBottom == 0 {
				symbol = unicodeTop
			} else {
				symbol = unicodeMinimalBodyFullWick
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom == i && rWickBottom == i:
			if qrWickTop == 0 {
				symbol = unicodeFlat
			} else if qrBodyTop == 0 {
				symbol = unicodeUpperWick
			} else if qrBodyBottom == 0 {
				symbol = unicodeHalfBodyTop
			} else if qrWickBottom == 0 {
				symbol = unicodeMinimalBodyBottomHalfWick
			} else {
				symbol = unicodeMinimalBodyFullWick
			}

		// TOP

		case rWickTop > i && rBodyTop == i && rBodyBottom == i && rWickBottom == i:
			if qrBodyTop == 0 {
				symbol = unicodeWick
			} else if qrBodyBottom == 0 {
				symbol = unicodeTop
			} else if qrWickBottom == 0 {
				symbol = unicodeMinimalBodyFullWick
			} else {
				symbol = unicodeMinimalBodyTopHalfWick
			}

		case rWickTop > i && rBodyTop == i && rBodyBottom < i && rWickBottom < i:
			if qrBodyTop == 0 {
				symbol = unicodeWick
			} else {
				symbol = unicodeTop
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom < i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = unicodeVoid
			} else if qrBodyTop == 0 {
				symbol = unicodeUpperWick
			} else {
				symbol = unicodeHalfBodyTop
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom == i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = unicodeVoid
			} else if qrBodyTop == 0 {
				symbol = unicodeUpperWick
			} else if qrBodyBottom == 0 {
				symbol = unicodeHalfBodyTop
			} else {
				symbol = unicodeMinimalBodyBottomHalfWick
			}

		case rWickTop == i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = unicodeVoid
			} else {
				symbol = unicodeUpperWick
			}

		case rWickTop > i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			symbol = unicodeWick

		default:
			symbol = "?"
		}

		symbols[int(i)] = symbol
	}

	return column{
		symbols: symbols,
		isUp:    c.Open <= c.Close,
	}
}
