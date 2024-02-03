package candlesticks

import (
	"github.com/fatih/color"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
)

type Chart struct {
	height, width int

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

func (chart *Chart) SetHeight(height int) {
	chart.height = height
}

func (chart *Chart) SetWidth(width int) {
	chart.width = width
}

func (chart Chart) Grid() charts.Grid {
	columns := chart.toColumns()

	grid := charts.NewGrid(chart.height, chart.width)
	for y := 0; y < chart.height; y++ {
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
	dataEnd := chart.cursor + chart.width
	if len(chart.data) < dataEnd {
		dataEnd = len(chart.data)
	}

	min, max := getMinMax(chart.data[dataStart:dataEnd])

	newData := make([]column, dataEnd)
	for i, c := range chart.data[dataStart:dataEnd] {
		newData[i] = newColumn(c, min, max, chart.height)
		if i == int(dataEnd-1) {
			break
		}
	}

	return newData
}

func (chart Chart) View() string {
	return chart.Grid().View()
}
