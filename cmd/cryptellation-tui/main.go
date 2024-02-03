package main

import (
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts/candlesticks"
)

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

type App struct {
	canvas     charts.Canvas
	csChart    candlesticks.Chart
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
	a.csChart = candlesticks.NewChart(candlesticks.ExampleData)
	a.canvas = charts.NewCanvas(time.Now().Add(-20*time.Minute), time.Now())

	return tea.ClearScreen
}

func (a *App) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left):
			a.csChart.MoveGridLeft()
		case key.Matches(msg, keys.Right):
			a.csChart.MoveGridRight()
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

	a.canvas.SetHeight(a.windowSize.Height - helpViewHeight)
	a.canvas.SetWidth(a.windowSize.Width)
	a.canvas.AddChart(&a.csChart)

	return a.canvas.View() + helpView
}
