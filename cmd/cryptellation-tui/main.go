package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

type model int

func main() {
	p := tea.NewProgram(model(5), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	c := toDiagram(globalData, 10, 20)
	return display(c)
}
