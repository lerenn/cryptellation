package main

import "github.com/charmbracelet/bubbles/key"

var (
	helpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)

	quitKey = key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	)
)

type keyMap struct {
	genericBindings  []key.Binding
	specificBindings []key.Binding
}

func newKeyMap(keys []key.Binding) keyMap {
	return keyMap{
		genericBindings:  []key.Binding{helpKey, quitKey},
		specificBindings: keys,
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return k.genericBindings
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.specificBindings, // first column
		k.genericBindings,  // second column
	}
}
