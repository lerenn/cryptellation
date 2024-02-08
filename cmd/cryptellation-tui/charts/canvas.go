package charts

import (
	"fmt"
	"math"
	"time"

	"github.com/dsnet/golib/unitconv"
)

type Canvas struct {
	height, width int

	start time.Time
	delta time.Duration

	charts []Chart
}

func NewCanvas(start time.Time, delta time.Duration) *Canvas {
	return &Canvas{
		start: start,
		delta: delta,
	}
}

const (
	unicodeVerticalAxis               = "│"
	unicodeVerticalAxisPointExtension = "╶"
	unicodeVerticalAxisLegend         = "┤"

	unicodeVerticalHorizontalAxisJointure = "┼"

	unicodeHorizontalAxis       = "─"
	unicodeHorizontalAxisLegend = "┬"
)

const (
	horizontalLegendSize = 1
	horizontalLegendGap  = 20
	horizontalAxisGap    = horizontalLegendSize

	verticalLegendSize = 7
	verticalLegendGap  = 4
	verticalAxisGap    = verticalLegendSize + 1
)

func (canvas Canvas) View() string {
	// Update subcharts and vertical min/max
	min, max := canvas.updateSubCharts()
	if min == math.MaxFloat64 { // If value is not really set
		min, max = 0, 0
	}

	// Create a new grid
	g := NewGrid(canvas.height, canvas.width)

	// Vertical axis line
	for i := horizontalAxisGap; i < canvas.height; i++ {
		g.SetSlotCharacterIfExists(verticalAxisGap, i, unicodeVerticalAxis)
	}

	// Horizontal axis line
	for i := verticalAxisGap; i < canvas.width; i++ {
		g.SetSlotCharacterIfExists(i, horizontalAxisGap, unicodeHorizontalAxis)
	}

	// Join between axis
	g.InsertCharacter(verticalAxisGap-1, horizontalAxisGap,
		unicodeVerticalAxisPointExtension,
		unicodeVerticalHorizontalAxisJointure)

	// Vertical legend
	// total := int(max-min) / canvas.height
	total := canvas.height / verticalLegendGap
	valueGap := (max - min) / float64(total)
	for i := 0; i < total; i++ {
		if i != 0 { // The first is already set as cross jointure
			g.InsertCharacter(verticalAxisGap-1, horizontalAxisGap+verticalLegendGap*i,
				unicodeVerticalAxisPointExtension, unicodeVerticalAxisLegend)
		}

		value := min + valueGap*float64(i)
		text := unitconv.FormatPrefix(value, unitconv.SI, 2)

		missingSpacesCount := verticalLegendSize - len(text)
		for i := 0; i < missingSpacesCount; i++ {
			text = " " + text
		}

		g.InsertText(0, horizontalAxisGap+verticalLegendGap*i, text)
	}

	// Horizontal legend
	total = canvas.width / horizontalLegendGap
	for i := 0; i < total; i++ {
		if i != 0 { // The first is already set as cross jointure
			g.InsertCharacter(verticalAxisGap+i*horizontalLegendGap, horizontalAxisGap, unicodeHorizontalAxisLegend)
		}

		t := canvas.start.Add(canvas.delta * time.Duration(i))
		month, day, hour, minute, second := t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()
		g.InsertText(verticalAxisGap+i*horizontalLegendGap, 0, fmt.Sprintf("%02d:%02d:%02d [%d/%d]", hour, minute, second, day, int(month)))
	}

	// Generate subcharts
	for _, c := range canvas.charts {
		g.ApplySubGrid(verticalAxisGap+1, horizontalLegendSize+1, c.Grid())
	}

	return g.View()
}

func (canvas *Canvas) AddChart(chart Chart) {
	canvas.charts = append(canvas.charts, chart)
	canvas.updateSubChartSize(len(canvas.charts) - 1)
}

func (canvas *Canvas) updateSubCharts() (min, max float64) {
	// Update start time
	for _, c := range canvas.charts {
		c.SetDisplayedTime(canvas.start)
	}

	// Get the minimal vertical data and the max vertical data
	min, max = math.MaxFloat64, -math.MaxFloat64
	for _, c := range canvas.charts {
		chartMin, chartMax := c.GetDisplayedDataMinMax()
		if chartMin < min {
			min = chartMin
		}

		if chartMax > max {
			max = chartMax
		}
	}

	// Update size and vertical boundaries
	for i, c := range canvas.charts {
		canvas.updateSubChartSize(i)
		c.SetVerticalBoundaries(min, max)
	}

	return
}

func (canvas *Canvas) updateSubChartSize(i int) {
	height := canvas.height - (horizontalLegendSize + 1)
	if height > 0 {
		canvas.charts[i].SetHeight(height)
	}

	width := canvas.width - (verticalAxisGap + 1)
	if width > 0 {
		canvas.charts[i].SetWidth(width)
	}
}

func (canvas *Canvas) SetHeight(height int) {
	canvas.height = height
}

func (canvas *Canvas) SetWidth(width int) {
	canvas.width = width
}

func (canvas *Canvas) MoveLeft() {
	for _, c := range canvas.charts {
		c.MoveLeft()
	}
	canvas.start = canvas.start.Add(-canvas.delta)
}

func (canvas *Canvas) MoveRight() {
	for _, c := range canvas.charts {
		c.MoveRight()
	}
	canvas.start = canvas.start.Add(canvas.delta)
}
