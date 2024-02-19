package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

type App struct {
	subview    View
	Program    *tea.Program
	windowSize tea.WindowSizeMsg
	help       help.Model
}

type dataUpdate struct{}

func main() {
	app := &App{}
	p := tea.NewProgram(app, tea.WithAltScreen())
	app.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Init() tea.Cmd {
	a.subview = &empty{}
	a.subview = NewCandlesticksView(a.Program)
	return tea.ClearScreen
}

func (a *App) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpKey):
			a.help.ShowAll = !a.help.ShowAll
		case key.Matches(msg, quitKey):
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.help.Width = msg.Width
		a.windowSize = msg
	}

	if a.subview != nil {
		a.subview.Update(message)
	}

	return a, nil
}

func (a *App) View() string {
	if a.windowSize.Height == 0 || a.windowSize.Width == 0 {
		return ""
	}

	// Generate help view
	helpView := a.help.View(newKeyMap(a.subview.Keys()))
	helpViewHeight := strings.Count(helpView, "\n") + 1

	subview := a.subview.View(0, helpViewHeight)

	return subview + helpView
}
