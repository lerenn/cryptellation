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
	if y < g.Height() && x < g.Width() {
		return &g[y][x]
	}
	return nil
}

func (g Grid) View() string {
	str := ""
	for y := g.Height() - 1; y >= 0; y-- {
		for x := 0; x < g.Width(); x++ {
			s := g.Slot(x, y)
			if s != nil {
				str += s.Character
			} else {
				str += ""
			}
		}
		str += "\n"
	}
	return str
}

func (grid Grid) Width() int {
	if len(grid) > 0 {
		return len(grid[0])
	}
	return 0
}

func (grid Grid) Height() int {
	return len(grid)
}

func (grid *Grid) InsertText(x, y int, text string) {
	for i := 0; i < len(text); i++ {
		grid.SetSlotCharacterIfExists(x+i, y, string(text[i]))
	}
}

func (grid *Grid) InsertCharacter(x, y int, chars ...string) {
	for i := 0; i < len(chars); i++ {
		grid.SetSlotCharacterIfExists(x+i, y, chars[i])
	}
}

func (grid *Grid) SetSlotCharacterIfExists(x, y int, character string) {
	s := grid.Slot(x, y)
	if s != nil {
		s.Character = character
	}
}

func (grid *Grid) ApplySubGrid(x, y int, g Grid) {
	for i := 0; i < g.Width() && i < grid.Width(); i++ {
		for j := 0; j < g.Height() && j < grid.Height(); j++ {
			s := g.Slot(i, j)
			if s != nil {
				grid.InsertCharacter(x+i, y+j, s.Character)
			}
		}
	}
}
