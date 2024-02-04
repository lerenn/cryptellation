package charts

import "time"

type Chart interface {
	SetHeight(int)
	SetWidth(int)
	Grid() Grid
	MoveRight()
	MoveLeft()
	SetVerticalBoundaries(min, max float64)
	GetDisplayedDataMinMax() (min, max float64)
	SetDisplayedTime(t time.Time)
}
