package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/views"
)

type App struct {
	Program    *tea.Program
	windowSize tea.WindowSizeMsg

	subview views.View
	cmdBar  CommandBar
}

func (a *App) Init() tea.Cmd {
	a.subview = &views.Help{}
	return tea.ClearScreen
}

func (a *App) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		a.cmdBar.Update(msg)
		if !a.cmdBar.Enabled() {
			a.subview.Update(message)
		}

	case Command:
		switch msg.Name {
		case "q", "quit":
			return a, tea.Quit
		case "c", "candlesticks":
			if len(msg.Arguments) == 0 {
				a.cmdBar.AddError("Missing exchange name")
			} else if len(msg.Arguments) == 1 {
				a.cmdBar.AddError("Missing pair")
			} else if len(msg.Arguments) == 2 {
				a.cmdBar.AddError("Missing period")
			} else {
				a.subview = views.NewCandlesticksView(
					a.Program,
					msg.Arguments[0],
					msg.Arguments[1],
					msg.Arguments[2])
				a.subview.Update(a.windowSize)
			}
		default:
			a.cmdBar.AddError(fmt.Sprintf("Unknown command: %q", msg.Name))
		}

	case tea.WindowSizeMsg:
		a.windowSize = msg
		a.subview.Update(message)
	}

	return a, nil
}

func (a *App) View() string {
	// Skip this step if there is no height or no width
	if a.windowSize.Height == 0 || a.windowSize.Width == 0 {
		return ""
	}

	// Generate command bar if required
	commandBarView := a.cmdBar.View(a.windowSize.Width)

	// Display subview
	usedHeight := strings.Count(commandBarView, "\n")
	subview := a.subview.View(0, usedHeight)

	return subview + commandBarView
}
