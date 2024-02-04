package charts

type Chart interface {
	SetHeight(int)
	SetWidth(int)
	Grid() Grid
	MoveRight()
	MoveLeft()
}
