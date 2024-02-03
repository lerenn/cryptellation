package candlesticks

import (
	"github.com/fatih/color"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
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

func (chart *Chart) MoveGridLeft() {
	if chart.cursor > 0 {
		chart.cursor--
	}
}

func (chart *Chart) MoveGridRight() {
	if chart.cursor < len(chart.data)-1 {
		chart.cursor++
	}
}

func (chart Chart) Grid() charts.Grid {
	columns := chart.toColumns()

	grid := charts.NewGrid(chart.Height, chart.Width)
	for y := 0; y < chart.Height; y++ {
		for x, c := range columns {
			// If the column is empty, doesn't display anything
			if len(c.symbols) == 0 {
				continue
			}

			if c.isUp {
				grid.InsertCharacter(x, y, color.GreenString(c.symbols[y]))
			} else {
				grid.InsertCharacter(x, y, color.RedString(c.symbols[y]))
			}
		}
	}
	return grid
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

func (chart Chart) View() string {
	return chart.Grid().View()
}
