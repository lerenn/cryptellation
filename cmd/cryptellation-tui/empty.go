package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type empty struct {
	windowSize tea.WindowSizeMsg
}

func (e empty) Keys() []key.Binding {
	return []key.Binding{}
}

func (e *empty) Update(message tea.Msg) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		e.windowSize = msg
	}
}

func (e empty) View(_, yPad int) string {
	return strings.Repeat("\n", e.windowSize.Height-yPad)
}
