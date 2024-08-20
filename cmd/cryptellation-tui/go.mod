module github.com/lerenn/cryptellation/cmd/cryptellation-tui

go 1.22.4

replace (
	github.com/lerenn/cryptellation/clients/go => ../../clients/go
	github.com/lerenn/cryptellation/pkg => ../../pkg

	github.com/lerenn/cryptellation/client => ../../svc/backtests
	github.com/lerenn/cryptellation/candlesticks => ../../svc/candlesticks
	github.com/lerenn/cryptellation/exchanges => ../../svc/exchanges
	github.com/lerenn/cryptellation/forwardtests => ../../svc/forwardtests
	github.com/lerenn/cryptellation/indicators => ../../svc/indicators
	github.com/lerenn/cryptellation/ticks => ../../svc/ticks
)

require (
	github.com/lerenn/cryptellation/pkg v0.0.0-00010101000000-000000000000
	github.com/lerenn/cryptellation/candlesticks v0.0.0-00010101000000-000000000000
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/bubbletea v0.26.6
	github.com/dsnet/golib/unitconv v1.0.2
	github.com/fatih/color v1.17.0
)

require (
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/charmbracelet/x/ansi v0.1.2 // indirect
	github.com/charmbracelet/x/input v0.1.0 // indirect
	github.com/charmbracelet/x/term v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.1.0 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)
