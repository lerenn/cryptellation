package charts

import (
	"fmt"
	"math"
	"time"
)

type Canvas struct {
	height, width int

	start time.Time
	delta time.Duration

	charts []Chart
}

func NewCanvas(start time.Time, delta time.Duration) Canvas {
	return Canvas{
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

	verticalLegendSize = 6
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
		g.InsertText(0, horizontalAxisGap+verticalLegendGap*i, fmt.Sprintf("%.1f", min+valueGap*float64(i)))
	}

	// Horizontal legend
	timeWidth := canvas.delta * time.Duration(canvas.width)
	total = canvas.width / horizontalLegendGap
	timeGap := timeWidth / time.Duration(total)
	for i := 0; i < total; i++ {
		h := canvas.start.Add(timeGap * time.Duration(i)).Hour()
		m := canvas.start.Add(timeGap * time.Duration(i)).Minute()
		s := canvas.start.Add(timeGap * time.Duration(i)).Second()
		if i != 0 { // The first is already set as cross jointure
			g.InsertCharacter(verticalAxisGap+i*horizontalLegendGap, horizontalAxisGap, unicodeHorizontalAxisLegend)
		}
		g.InsertText(verticalAxisGap+i*horizontalLegendGap, 0, fmt.Sprintf("%02d:%02d:%02d", h, m, s))
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
