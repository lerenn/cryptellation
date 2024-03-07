package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Empty struct {
	windowSize tea.WindowSizeMsg
}

func (e Empty) Keys() []key.Binding {
	return []key.Binding{}
}

func (e *Empty) Update(message tea.Msg) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		e.windowSize = msg
	}
}

func (e Empty) View(_, yPad int) string {
	if e.windowSize.Height < yPad {
		return ""
	}
	return strings.Repeat("\n", e.windowSize.Height-yPad)
}
