package views

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/utils"

	cdsclient "github.com/lerenn/cryptellation/candlesticks/clients/go"
	candlestickscache "github.com/lerenn/cryptellation/candlesticks/clients/go/cache"
	candlesticksnats "github.com/lerenn/cryptellation/candlesticks/clients/go/nats"
	candlesticksretry "github.com/lerenn/cryptellation/candlesticks/clients/go/retry"
	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/candlesticks/pkg/period"

	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts"
	"github.com/lerenn/cryptellation/cmd/cryptellation-tui/charts/candlesticks"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type candlesticksDataUpdate struct{}

type CandlesticksView struct {
	client           cdsclient.Client
	updateInProgress bool
	canvas           *charts.Canvas
	chart            *candlesticks.Chart
	windowSize       tea.WindowSizeMsg

	exchange string
	pair     string
	period   period.Symbol

	program *tea.Program
}

func NewCandlesticksView(program *tea.Program, exchange, pair, periodSymbol string) *CandlesticksView {
	client, err := candlesticksnats.New(config.LoadNATS())
	if err != nil {
		log.Fatal(err)
	}
	client = candlestickscache.New(client)
	client = candlesticksretry.New(client)

	per := period.Symbol(strings.ToUpper(periodSymbol))
	if err := per.Validate(); err != nil {
		log.Fatal(err)
	}

	cv := &CandlesticksView{
		client:   client,
		program:  program,
		exchange: strings.ToLower(exchange),
		pair:     strings.ToUpper(pair),
		period:   per,
	}

	cv.chart = candlesticks.NewChart(&candlestick.List{}, per)

	cv.canvas = charts.NewCanvas(utils.Must(time.Parse(time.RFC3339, "2024-01-01T01:00:00Z")), per.Duration())
	cv.canvas.AddChart(cv.chart)

	return cv
}

func (cv *CandlesticksView) moveCount() int {
	return cv.windowSize.Width * 2 / 3
}

func (cv *CandlesticksView) Update(message tea.Msg) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, candlestickKeyLeft):
			for i := 0; i < cv.moveCount(); i++ {
				cv.canvas.MoveLeft()
			}
		case key.Matches(msg, candlestickKeyRight):
			for i := 0; i < cv.moveCount(); i++ {
				cv.canvas.MoveRight()
			}
		}

	case tea.WindowSizeMsg:
		cv.windowSize = msg
	}

	cv.updateMissingCandlesticks()
}

func (cv *CandlesticksView) updateMissingCandlesticks() {
	first, last := cv.chart.MissingData(cv.windowSize.Width)
	if first != nil && last != nil {
		go func() {
			if cv.updateInProgress {
				return
			}
			cv.updateInProgress = true
			defer func() { cv.updateInProgress = false }()

			delta := time.Duration(cv.windowSize.Width)
			first = utils.ToReference(first.Add(-cv.period.Duration() * delta))
			last = utils.ToReference(last.Add(cv.period.Duration() * delta))

			list, err := cv.client.Read(context.TODO(), cdsclient.ReadCandlesticksPayload{
				Exchange: cv.exchange,
				Pair:     cv.pair,
				Period:   cv.period,
				Start:    first,
				End:      last,
			})
			if err != nil {
				return
			}

			if err := cv.chart.UpsertData(list); err != nil {
				log.Fatal(err)
			}

			// Send the main program an update
			cv.program.Send(candlesticksDataUpdate{})
		}()
	}
}

func (cv *CandlesticksView) View(xPad, yPad int) string {
	cv.canvas.SetHeight(cv.windowSize.Height - yPad)
	cv.canvas.SetWidth(cv.windowSize.Width - xPad)

	return cv.canvas.View()
}

var (
	candlestickKeyLeft = key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	)

	candlestickKeyRight = key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	)
)

func (cv CandlesticksView) Keys() []key.Binding {
	return []key.Binding{
		candlestickKeyLeft, candlestickKeyRight,
	}
}
