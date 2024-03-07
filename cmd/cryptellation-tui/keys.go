package main

import "github.com/charmbracelet/bubbles/key"

var (
	backspaceKey = key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("<backspace>", "remove character"),
	)

	commandKey = key.NewBinding(
		key.WithKeys(":"),
		key.WithHelp("<:>", "activate command bar"),
	)

	enterKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "execute command"),
	)

	escapeKey = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<escape>", "close command bar"),
	)
)
