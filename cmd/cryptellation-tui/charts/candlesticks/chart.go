package candlesticks

import (
	"time"

	"github.com/fatih/color"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type Chart struct {
	height, width            int
	verticalMin, verticalMax float64

	Data   *candlestick.List
	cursor time.Time
	period period.Symbol
}

func NewChart(data *candlestick.List, period period.Symbol) *Chart {
	return &Chart{
		Data:   data,
		period: period,
	}
}

func (chart *Chart) UpsertData(data *candlestick.List) error {
	if chart.Data.Len() == 0 {
		chart.Data = data
		return nil
	}

	return chart.Data.Merge(data, nil)
}

func (chart *Chart) MoveLeft() {
	chart.cursor = chart.cursor.Add(-chart.period.Duration())
}

func (chart *Chart) MoveRight() {
	chart.cursor = chart.cursor.Add(chart.period.Duration())
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
	start, end := chart.displayedStartEnd()
	newData := make([]column, chart.width)

	for current, i := start, 0; current.Before(end); current, i = current.Add(chart.period.Duration()), i+1 {
		c, exists := chart.Data.Get(current)
		if exists {
			newData[i] = newColumn(c, chart.verticalMin, chart.verticalMax, chart.height)
		}
	}

	return newData
}

func (chart Chart) displayedStartEnd() (start, end time.Time) {
	start = chart.cursor
	end = start.Add(chart.period.Duration() * time.Duration(chart.width))
	return
}

func (chart Chart) MissingData() (first, last *time.Time) {
	start, end := chart.displayedStartEnd()

	for current := start; current.Before(end); current = current.Add(chart.period.Duration()) {
		_, exists := chart.Data.Get(current)
		if !exists {
			copyCurrent := current
			if first == nil {
				first = &copyCurrent
			}
			last = &copyCurrent
		}
	}

	return
}

func (chart Chart) GetDisplayedDataMinMax() (min, max float64) {
	start, end := chart.displayedStartEnd()

	data := candlestick.NewEmptyListFrom(chart.Data)
	for current := start; current.Before(end); current = current.Add(chart.period.Duration()) {
		c, exists := chart.Data.Get(current)
		if exists {
			data.Set(current, c)
		}
	}

	return getMinMax(data)
}

func (chart *Chart) SetVerticalBoundaries(min, max float64) {
	chart.verticalMin = min
	chart.verticalMax = max
}

func (chart *Chart) SetDisplayedTime(t time.Time) {
	chart.cursor = t
}

func (chart Chart) View() string {
	return chart.Grid().View()
}
