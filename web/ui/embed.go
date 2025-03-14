package ui

import (
	"embed"
)

// StaticFS is the embedded file system for the UI React App.
//
//go:embed build
var StaticFS embed.FS
