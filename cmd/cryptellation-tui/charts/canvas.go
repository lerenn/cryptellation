package charts

import (
	"fmt"
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
	unicodeVerticalAxis                   = "│"
	unicodeVerticalAxisPointExtension     = "╶"
	unicodeVerticalHorizontalAxisJointure = "┼"
	unicodeHorizontalAxis                 = "─"
	unicodeHorizontalAxisLegend           = "┬"
)

const (
	horizontalLegendSize = 1
	horizontalLegendGap  = 20
	horizontalAxisGap    = horizontalLegendSize

	verticalLegendSize = 6
	verticalAxisGap    = verticalLegendSize + 1
)

func (canvas Canvas) View() string {
	// Update subcharts
	for i := range canvas.charts {
		canvas.updateSubChartSize(i)
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
	g.InsertText(0, horizontalAxisGap, "000.00")

	// Horizontal legend
	timeWidth := canvas.delta * time.Duration(canvas.width)
	total := canvas.width / horizontalLegendGap
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

	for _, c := range canvas.charts {
		g.ApplySubGrid(verticalAxisGap+1, horizontalLegendSize+1, c.Grid())
	}

	return g.View()
}

func (canvas *Canvas) AddChart(chart Chart) {
	canvas.charts = append(canvas.charts, chart)
	canvas.updateSubChartSize(len(canvas.charts) - 1)
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
