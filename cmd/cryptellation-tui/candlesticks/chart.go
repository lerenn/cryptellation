package candlesticks

import (
	"github.com/fatih/color"
)

type Chart struct {
	Height, Width int

	data   []Candlestick
	cursor int
}

func NewChart(data []Candlestick) Chart {
	return Chart{
		data: data,
	}
}

func (chart *Chart) MoveViewLeft() {
	if chart.cursor > 0 {
		chart.cursor--
	}
}

func (chart *Chart) MoveViewRight() {
	if chart.cursor < len(chart.data)-1 {
		chart.cursor++
	}
}

func (chart Chart) View() string {
	columns := chart.toColumns()

	str := ""
	for i := chart.Height - 1; i >= 0; i-- {
		for _, c := range columns {
			// If the column is empty, doesn't display anything
			if len(c.symbols) == 0 {
				str += unicodeVoid
				continue
			}

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

func (chart Chart) toColumns() []column {
	dataStart := chart.cursor
	dataEnd := chart.cursor + chart.Width
	if len(chart.data) < dataEnd {
		dataEnd = len(chart.data)
	}

	min, max := getMinMax(chart.data[dataStart:dataEnd])

	newData := make([]column, dataEnd)
	for i, c := range chart.data[dataStart:dataEnd] {
		newData[i] = newColumn(c, min, max, chart.Height)
		if i == int(dataEnd-1) {
			break
		}
	}

	return newData
}
