package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := &App{}

	// Add command line
	app.cmdBar = NewCommandBar(app)

	// Add buubletea program
	app.Program = tea.NewProgram(app, tea.WithAltScreen())
	if _, err := app.Program.Run(); err != nil {
		log.Fatal(err)
	}
}
