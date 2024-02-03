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
	g := NewGrid(canvas.Height, canvas.Width)

	// Vertical axis line
	for i := 0; i < canvas.Height-horizontalAxisGap; i++ {
		g.SetSlotCharacterIfExists(verticalAxisGap, i, unicodeVerticalAxis)
	}

	// Horizontal axis line
	lastVerticalRowNb := canvas.Height - (1 + horizontalAxisGap)
	for i := verticalAxisGap - 1; i < canvas.Width; i++ {
		g.SetSlotCharacterIfExists(i, lastVerticalRowNb, unicodeHorizontalAxis)
	}

	// Join between axis
	g.InsertCharacter(verticalAxisGap-1, lastVerticalRowNb-horizontalAxisGap+1,
		unicodeVerticalAxisPointExtension,
		unicodeVerticalHorizontalAxisJointure)

	// Vertical legend
	g.InsertText(0, canvas.Height-2, "000.00")

	// Horizontal legend
	timeWidth := canvas.end.Sub(canvas.start)
	total := canvas.Width / horizontalLegendGap
	timeGap := timeWidth / time.Duration(total)
	for i := 0; i < total; i++ {
		h := canvas.start.Add(timeGap * time.Duration(i)).Hour()
		m := canvas.start.Add(timeGap * time.Duration(i)).Minute()
		s := canvas.start.Add(timeGap * time.Duration(i)).Second()
		if i != 0 { // The first is already set as cross jointure
			g.InsertCharacter(verticalAxisGap+i*horizontalLegendGap, canvas.Height-2, unicodeHorizontalAxisLegend)
		}
		g.InsertText(verticalAxisGap+i*horizontalLegendGap, canvas.Height-1, fmt.Sprintf("%02d:%02d:%02d", h, m, s))
	}

	return g.View()
}
