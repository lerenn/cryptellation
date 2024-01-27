package charts

type GridSlot struct {
	Character string
}

type Grid [][]GridSlot

func NewGrid(height, width int) Grid {
	g := make([][]GridSlot, height)
	for i := range g {
		g[i] = make([]GridSlot, width)
		for j := range g[i] {
			g[i][j].Character = " "
		}
	}
	return g
}

func (g Grid) Slot(x, y int) *GridSlot {
	return &g[x][y]
}

func (g Grid) View() string {
	str := ""
	for _, row := range g {
		for _, slot := range row {
			str += slot.Character
		}
		str += "\n"
	}
	return str
}

func (grid *Grid) InsertText(x, y int, text string) {
	for i := 0; i < len(text); i++ {
		grid.Slot(x, y+i).Character = string(text[i])
	}
}

func (grid *Grid) InsertCharactersHorizontally(x, y int, chars ...string) {
	for i := 0; i < len(chars); i++ {
		grid.Slot(x, y+i).Character = chars[i]
	}
}
