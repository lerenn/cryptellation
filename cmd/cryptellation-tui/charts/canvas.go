package charts

import (
	"fmt"
	"time"
)

type Canvas struct {
	height, width int

	start time.Time
	end   time.Time

	charts []Chart
}

func NewCanvas(start, end time.Time) Canvas {
	return Canvas{
		start: start,
		end:   end,
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
	timeWidth := canvas.end.Sub(canvas.start)
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
	canvas.charts[i].SetHeight(canvas.height - (horizontalLegendSize + 1))
	canvas.charts[i].SetWidth(canvas.width - (verticalAxisGap + 1))
}

func (canvas *Canvas) SetHeight(height int) {
	for i := range canvas.charts {
		canvas.updateSubChartSize(i)
	}
	canvas.height = height
}

func (canvas *Canvas) SetWidth(width int) {
	for i := range canvas.charts {
		canvas.updateSubChartSize(i)
	}
	canvas.width = width
}
