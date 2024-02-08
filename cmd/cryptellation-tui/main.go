package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts/candlesticks"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksclient "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

type App struct {
	CandlesticksClient           candlesticksclient.Client
	candlesticksUpdateInProgress bool

	canvas       *charts.Canvas
	candlesticks *candlesticks.Chart

	windowSize tea.WindowSizeMsg
	help       help.Model

	Program *tea.Program
}

type dataUpdate struct{}

func main() {
	candlesticksClient, err := candlesticksclient.NewClient(config.LoadNATS())
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		CandlesticksClient: candlesticksClient,
	}
	p := tea.NewProgram(app, tea.WithAltScreen())
	app.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Init() tea.Cmd {
	a.canvas = charts.NewCanvas(utils.Must(time.Parse(time.RFC3339, "2022-12-01T01:00:00Z")), time.Hour)

	a.candlesticks = candlesticks.NewChart(&candlestick.List{}, period.H1)
	a.canvas.AddChart(a.candlesticks)
	defer a.updateMissingCandlesticks()

	return tea.ClearScreen
}

func (a *App) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Left):
			a.canvas.MoveLeft()
		case key.Matches(msg, keys.Right):
			a.canvas.MoveRight()
		case key.Matches(msg, keys.Help):
			a.help.ShowAll = !a.help.ShowAll
		case key.Matches(msg, keys.Quit):
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.windowSize = msg
		a.help.Width = msg.Width
	}

	a.updateMissingCandlesticks()

	return a, nil
}

func (a *App) updateMissingCandlesticks() {
	first, last := a.candlesticks.MissingData(a.windowSize.Width)
	if first != nil && last != nil {
		go func() {
			if a.candlesticksUpdateInProgress {
				return
			}
			a.candlesticksUpdateInProgress = true
			defer func() { a.candlesticksUpdateInProgress = false }()

			delta := time.Duration(a.windowSize.Width)
			first = utils.ToReference(first.Add(-time.Hour * delta))
			last = utils.ToReference(last.Add(time.Hour * delta))

			list, err := a.CandlesticksClient.Read(context.TODO(), client.ReadCandlesticksPayload{
				Exchange: "binance",
				Pair:     "ETH-USDT",
				Period:   period.H1,
				Start:    first,
				End:      last,
			})
			if err != nil {
				return
			}
			a.candlesticks.UpsertData(list)
			a.Program.Send(dataUpdate{})
		}()
	}
}

func (a *App) View() string {
	if a.windowSize.Height == 0 || a.windowSize.Width == 0 {
		return ""
	}

	// Generate help view
	helpView := a.help.View(keys)
	helpViewHeight := strings.Count(helpView, "\n") + 1

	a.canvas.SetHeight(a.windowSize.Height - helpViewHeight)
	a.canvas.SetWidth(a.windowSize.Width)

	return a.canvas.View() + helpView
}
