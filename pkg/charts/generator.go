package charts

import (
	"context"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	client "github.com/lerenn/cryptellation/clients/go"
	cdksclient "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	indclient "github.com/lerenn/cryptellation/svc/indicators/clients/go"
)

type Generator struct {
	client client.Services
}

func NewGenerator(client client.Services) *Generator {
	return &Generator{
		client: client,
	}
}

type CandlesticksPayload struct {
	Name         string
	Candlesticks cdksclient.ReadCandlesticksPayload
}

func (c Generator) Candlesticks(ctx context.Context, payload CandlesticksPayload) (*charts.Kline, error) {
	l, err := c.client.Candlesticks().Read(ctx, payload.Candlesticks)
	if err != nil {
		return nil, err
	}

	x := make([]string, 0, l.Len())
	y := make([]opts.KlineData, 0, l.Len())
	if err := l.Loop(func(t time.Time, c candlestick.Candlestick) (bool, error) {
		x = append(x, t.Format("2006-01-02 15:04:05"))
		y = append(y, opts.KlineData{
			Value: [4]float64{c.Open, c.Close, c.Low, c.High},
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	chart := charts.NewKLine()
	chart.SetXAxis(x).AddSeries(payload.Name, y).SetSeriesOptions(
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "#60AF67",
			Color0:       "#DC6157",
			BorderColor:  "#60AF67",
			BorderColor0: "#DC6157",
		}),
	)
	chart.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	return chart, nil
}

type SMAPayload struct {
	SMA   indclient.SMAPayload
	Color string
}

func (c Generator) SMA(ctx context.Context, payload SMAPayload) (*charts.Line, error) {
	l, err := c.client.Indicators().SMA(ctx, payload.SMA)
	if err != nil {
		return nil, err
	}

	x := make([]string, 0, l.Len())
	y := make([]opts.LineData, 0, l.Len())
	if err := l.Loop(func(t time.Time, v float64) (bool, error) {
		x = append(x, t.Format("2006-01-02 15:04:05"))
		y = append(y, opts.LineData{
			Value: v,
		})
		return false, nil
	}); err != nil {
		return nil, err
	}

	chart := charts.NewLine()
	chart.SetXAxis(x).AddSeries("SMA", y).SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			Color: payload.Color,
		}))
	chart.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	return chart, nil
}
