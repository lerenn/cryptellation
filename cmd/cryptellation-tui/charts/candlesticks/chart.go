package candlesticks

import (
	"time"

	"cryptellation/svc/candlesticks/pkg/candlestick"
	"cryptellation/svc/candlesticks/pkg/period"

	"cryptellation/cmd/cryptellation-tui/charts"

	"github.com/fatih/color"
)

type Chart struct {
	height, width            int
	verticalMin, verticalMax float64

	data  *candlestick.List
	cache map[minMax]map[int64]column

	cursor time.Time
	period period.Symbol
}

type minMax struct {
	Min float64
	Max float64
}

func NewChart(data *candlestick.List, period period.Symbol) *Chart {
	return &Chart{
		data:   data,
		cache:  make(map[minMax]map[int64]column),
		period: period,
	}
}

func (chart *Chart) UpsertData(data *candlestick.List) error {
	if chart.data.Len() == 0 {
		chart.data = data
		return nil
	}

	return chart.data.Merge(data, nil)
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

func (chart *Chart) Grid() charts.Grid {
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

func (chart *Chart) toColumns() []column {
	start, end := chart.displayedStartEnd()
	min, max := chart.GetDisplayedDataMinMax()
	minMax := minMax{Min: min, Max: max}
	newData := make([]column, chart.width)

	if _, exists := chart.cache[minMax]; !exists {
		chart.cache[minMax] = make(map[int64]column, chart.width)
	}

	for current, i := start, 0; current.Before(end); current, i = current.Add(chart.period.Duration()), i+1 {
		if col, exists := chart.cache[minMax][current.Unix()]; exists {
			newData[i] = col
			continue
		}

		c, exists := chart.data.Get(current)
		if exists && c.High != 0 {
			col := newColumn(c, chart.verticalMin, chart.verticalMax, chart.height)
			newData[i], chart.cache[minMax][current.Unix()] = col, col
		}
	}

	return newData
}

func (chart Chart) displayedStartEnd() (start, end time.Time) {
	start = chart.cursor
	end = start.Add(chart.period.Duration() * time.Duration(chart.width))
	return
}

func (chart Chart) MissingData(margin int) (first, last *time.Time) {
	start, end := chart.displayedStartEnd()

	marginDuration := chart.period.Duration() * time.Duration(margin)
	start = start.Add(-marginDuration)
	end = end.Add(marginDuration)

	for current := start; current.Before(end); current = current.Add(chart.period.Duration()) {
		_, exists := chart.data.Get(current)
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

	data := candlestick.NewListFrom(chart.data)
	for current := start; current.Before(end); current = current.Add(chart.period.Duration()) {
		c, exists := chart.data.Get(current)
		if exists {
			_ = data.Set(current, c)
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
