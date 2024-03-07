package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Command struct {
	Name      string
	Arguments []string
}

type CommandBar struct {
	app          *App
	enabled      bool
	input        string
	errorMessage string
}

func NewCommandBar(app *App) CommandBar {
	return CommandBar{
		app: app,
	}
}

func (cb *CommandBar) Update(msg tea.KeyMsg) {
	// If not enabled, then wait for enabling and return
	if !cb.enabled {
		if key.Matches(msg, commandKey) {
			cb.enabled = true
		}
		return
	}

	// If enabled
	if key.Matches(msg, escapeKey) {
		cb.Disable()
	} else if key.Matches(msg, enterKey) {
		elements := strings.Split(cb.input, " ")
		cmd := Command{
			Name: elements[0],
		}

		if len(elements) > 0 {
			cmd.Arguments = elements[1:]
		}

		go cb.app.Program.Send(cmd)
		cb.Disable()
	} else if key.Matches(msg, backspaceKey) {
		if len(cb.input) > 0 {
			cb.input = cb.input[:len(cb.input)-1]
		}
	} else {
		cb.input += msg.String()
	}
}

func (cb *CommandBar) AddError(text string) {
	cb.errorMessage = text
}

func (cb *CommandBar) View(width int) string {
	var view string

	if cb.enabled {
		view += "> " + cb.input
	}

	if len(cb.errorMessage) > 0 {
		lenMiddle := width - len(cb.errorMessage) - len(view)
		view += strings.Repeat(" ", lenMiddle) + cb.errorMessage
	}

	return view
}

func (cb *CommandBar) Clear() {
	cb.input = ""
	cb.errorMessage = ""
}

func (cb CommandBar) Enabled() bool {
	return cb.enabled
}

func (cb *CommandBar) Disable() {
	cb.Clear()
	cb.enabled = false
}
