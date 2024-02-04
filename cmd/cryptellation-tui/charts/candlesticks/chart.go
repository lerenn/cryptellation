package candlesticks

import (
	"time"

	"github.com/fatih/color"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
)

type Chart struct {
	height, width            int
	verticalMin, verticalMax float64

	data   []*Candlestick
	cursor int
}

func NewChart(data []*Candlestick) Chart {
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
	data := chart.getDisplayedData()
	newData := make([]column, len(chart.data))
	for i, c := range data {
		if c == nil {
			continue
		}

		newData[i] = newColumn(*c, chart.verticalMin, chart.verticalMax, chart.height)
	}

	return newData
}

func (chart Chart) GetDisplayedDataMinMax() (min, max float64) {
	data := chart.getDisplayedData()
	return getMinMax(data)
}

func (chart *Chart) SetVerticalBoundaries(min, max float64) {
	chart.verticalMin = min
	chart.verticalMax = max
}

func (chart *Chart) SetDisplayedTime(t time.Time) {
	if len(chart.data) == 0 {
		return
	}

	delta := chart.data[0].Time.Sub(t)
	chart.cursor = -int(delta / time.Hour)
}

func (chart Chart) getDisplayedData() []*Candlestick {
	// Set start
	start := chart.cursor

	// Check if there is empty data before
	startGap := 0
	if start < 0 {
		startGap = -start
		start = 0
	}
	emptyStart := make([]*Candlestick, startGap)

	// Set end
	end := start + chart.width

	// Check if the end is after the end of the screen
	if len(chart.data) < end {
		end = len(chart.data)
	}
	end -= startGap // Remove the potential empty data gap before

	// Return empty data if the end if before start
	if end <= start {
		return []*Candlestick{}
	}

	return append(emptyStart, chart.data[start:end]...)
}

func (chart Chart) View() string {
	return chart.Grid().View()
}
