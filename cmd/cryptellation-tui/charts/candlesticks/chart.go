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

func (chart *Chart) MoveLeft() {
	chart.cursor--
}

func (chart *Chart) MoveRight() {
	chart.cursor++
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
	start := chart.cursor
	startGap := 0
	if start < 0 {
		startGap = -start
		start = 0
	}

	end := start + chart.width
	if len(chart.data) < end {
		end = len(chart.data)
	}

	newData := make([]column, end)
	if end <= start || startGap >= end {
		return newData
	}

	min, max := getMinMax(chart.data[start : end-startGap])
	for i, c := range chart.data[start : end-startGap] {
		newData[i+startGap] = newColumn(c, min, max, chart.height)
		if i == int(end-1) {
			break
		}
	}

	return newData
}

func (chart Chart) View() string {
	return chart.Grid().View()
}
