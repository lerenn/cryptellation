package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Help struct {
	windowSize tea.WindowSizeMsg
}

func (h Help) Keys() []key.Binding {
	return []key.Binding{}
}

func (h *Help) Update(message tea.Msg) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		h.windowSize = msg
	}
}

func (h Help) View(_, yPad int) string {
	if h.windowSize.Height < yPad {
		return ""
	}

	msg := `
 Cryptellation TUI
--------------------------------------------------------------------------------

Please type <:> to execute a command (<esc> if you want to quit command bar).
Then you can type one of the following command:

> candlesticks <exchange> <pair> <period>
  example: candlesticks binance eth-usdt h1 
	`

	msgHeight := strings.Count(msg, "\n")

	bottomMargin := strings.Repeat("\n", h.windowSize.Height-yPad-msgHeight)
	return msg + bottomMargin
}
