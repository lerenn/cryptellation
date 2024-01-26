package main

import (
	"math"
	"time"

	"github.com/fatih/color"
)

const (
	UnicodeVoid           = " "
	UnicodeBody           = "┃"
	UnicodeHalfBodyTop    = "╻"
	UnicodeHalfBodyBottom = "╹"
	UnicodeWick           = "│"
	UnicodeTop            = "╽"
	UnicodeBottom         = "╿"
	UnicodeUpperWick      = "╷"
	UnicodeLowerWick      = "╵"

	UnicodeFlat                      = "-"
	UnicodeMinimalBodyFullWick       = "┿"
	UnicodeMinimalBodyTopHalfWick    = "┷"
	UnicodeMinimalBodyBottomHalfWick = "┯"
)

type candlestick struct {
	Time  time.Time
	Open  float64
	High  float64
	Low   float64
	Close float64
}

func normalizePoint(point, min, max float64, height uint) float64 {
	point = (point - min) / (max - min)
	point = point * float64(height)
	return point
}

type column struct {
	symbols []string
	isUp    bool
}

func newColumn(c candlestick, min, max float64, height uint) column {
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
	for i := float64(0); i < float64(height); i++ {
		var symbol string

		switch {
		case rWickTop < i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			symbol = UnicodeVoid
		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom > i:
			symbol = UnicodeVoid

		// BOTTOM

		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom == i:
			if qrWickBottom-i == 0 {
				symbol = UnicodeWick
			} else {
				symbol = UnicodeLowerWick
			}

		case rWickTop > i && rBodyTop > i && rBodyBottom > i && rWickBottom < i:
			symbol = UnicodeWick

		case rWickTop > i && rBodyTop > i && rBodyBottom == i && rWickBottom == i:
			if qrBodyBottom == 0 {
				symbol = UnicodeBody
			} else if qrWickBottom == 0 {
				symbol = UnicodeBottom
			} else {
				symbol = UnicodeHalfBodyBottom
			}

		// MIDDLE

		case rWickTop > i && rBodyTop > i && rBodyBottom == i && rWickBottom < i:
			if qrBodyBottom == 0 {
				symbol = UnicodeBody
			} else {
				symbol = UnicodeBottom
			}

		case rWickTop > i && rBodyTop > i && rBodyBottom < i && rWickBottom < i:
			symbol = UnicodeBody

		case rWickTop > i && rBodyTop == i && rBodyBottom == i && rWickBottom < i:
			if qrBodyTop == 0 {
				symbol = UnicodeWick
			} else if qrBodyBottom == 0 {
				symbol = UnicodeTop
			} else {
				symbol = UnicodeMinimalBodyFullWick
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom == i && rWickBottom == i:
			if qrWickTop == 0 {
				symbol = UnicodeFlat
			} else if qrBodyTop == 0 {
				symbol = UnicodeUpperWick
			} else if qrBodyBottom == 0 {
				symbol = UnicodeHalfBodyTop
			} else if qrWickBottom == 0 {
				symbol = UnicodeMinimalBodyBottomHalfWick
			} else {
				symbol = UnicodeMinimalBodyFullWick
			}

		// TOP

		case rWickTop > i && rBodyTop == i && rBodyBottom == i && rWickBottom == i:
			if qrBodyTop == 0 {
				symbol = UnicodeWick
			} else if qrBodyBottom == 0 {
				symbol = UnicodeTop
			} else if qrWickBottom == 0 {
				symbol = UnicodeMinimalBodyFullWick
			} else {
				symbol = UnicodeMinimalBodyTopHalfWick
			}

		case rWickTop > i && rBodyTop == i && rBodyBottom < i && rWickBottom < i:
			if qrBodyTop == 0 {
				symbol = UnicodeWick
			} else {
				symbol = UnicodeTop
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom < i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = UnicodeVoid
			} else if qrBodyTop == 0 {
				symbol = UnicodeUpperWick
			} else {
				symbol = UnicodeHalfBodyTop
			}

		case rWickTop == i && rBodyTop == i && rBodyBottom == i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = UnicodeVoid
			} else if qrBodyTop == 0 {
				symbol = UnicodeUpperWick
			} else if qrBodyBottom == 0 {
				symbol = UnicodeHalfBodyTop
			} else {
				symbol = UnicodeMinimalBodyBottomHalfWick
			}

		case rWickTop == i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			if qrWickTop == 0 {
				symbol = UnicodeVoid
			} else {
				symbol = UnicodeUpperWick
			}

		case rWickTop > i && rBodyTop < i && rBodyBottom < i && rWickBottom < i:
			symbol = UnicodeWick

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

func toDiagram(data []candlestick, height, width uint) []column {
	if len(data) < int(width) {
		width = uint(len(data))
	}

	min, max := getMinMax(data[:width])

	newData := make([]column, width)
	for i, c := range data {
		newData[i] = newColumn(c, min, max, height)
		if i == int(width-1) {
			break
		}
	}

	return newData
}

func display(columns []column) string {
	str := ""
	for i := len(columns[0].symbols) - 1; i >= 0; i-- {
		for _, c := range columns {
			if c.isUp {
				str += color.GreenString(c.symbols[i])
			} else {
				str += color.RedString(c.symbols[i])
			}
		}
		str += "\n"
	}
	return str
}

func getMinMax(data []candlestick) (min float64, max float64) {
	min, max = math.MaxFloat64, 0
	for _, d := range data {
		if d.Low < min {
			min = d.Low
		}
		if d.High > max {
			max = d.High
		}
	}
	return
}
