package views

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Keys() []key.Binding
	Update(message tea.Msg)
	View(xPad, yPad int) string
}
