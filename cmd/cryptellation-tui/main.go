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
	cursor     int
	windowSize tea.WindowSizeMsg
	help       help.Model
}

func main() {
	p := tea.NewProgram(&App{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Init() tea.Cmd {
	return tea.ClearScreen
}

func (a *App) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left):
			if a.cursor > 0 {
				a.cursor--
			}
		case key.Matches(msg, keys.Right):
			if a.cursor < len(globalData)-1 {
				a.cursor++
			}
		case key.Matches(msg, keys.Help):
			a.help.ShowAll = !a.help.ShowAll
		case key.Matches(msg, keys.Quit):
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.windowSize = msg
		a.help.Width = msg.Width
	}

	return a, nil
}

func (a *App) View() string {
	if a.windowSize.Height == 0 {
		return ""
	}

	// Generate help view
	helpView := a.help.View(keys)
	helpViewHeight := strings.Count(helpView, "\n") + 1

	c := toDiagram(globalData[a.cursor:], uint(a.windowSize.Height-helpViewHeight), uint(a.windowSize.Width))
	return display(c) + helpView
}
