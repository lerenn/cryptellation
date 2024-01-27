package charts

import (
	"fmt"
	"time"
)

type Canvas struct {
	Height, Width int

	start time.Time
	end   time.Time
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
)

const (
	horizontalLegendSize = 1
	horizontalAxisGap    = horizontalLegendSize

	verticalLegendSize = 6
	verticalAxisGap    = verticalLegendSize + 1
)

func (canvas Canvas) View() string {
	g := NewGrid(canvas.Height, canvas.Width)

	// Vertical axis line
	for i := 0; i < canvas.Height-horizontalAxisGap; i++ {
		g.Slot(i, verticalAxisGap).Character = unicodeVerticalAxis
	}

	// Horizontal axis line
	lastVerticalRowNb := canvas.Height - (1 + horizontalAxisGap)
	for i := verticalAxisGap - 1; i < canvas.Width; i++ {
		g.Slot(lastVerticalRowNb, i).Character = unicodeHorizontalAxis
	}

	// Join between axis
	g.InsertCharactersHorizontally(lastVerticalRowNb, verticalAxisGap-1,
		unicodeVerticalAxisPointExtension,
		unicodeVerticalHorizontalAxisJointure)

	// Vertical legend
	g.InsertText(canvas.Height-2, 0, "000.00")

	// Horizontal legend
	g.InsertText(canvas.Height-1, verticalAxisGap,
		fmt.Sprintf("%02d:%02d:%02d", canvas.start.Hour(), canvas.start.Minute(), canvas.start.Second()))

	return g.View()
}
