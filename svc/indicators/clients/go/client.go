// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=client.go -destination=mock.gen.go -package client

package client

import (
	"context"
	"time"

	client "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

type Client interface {
	SMA(ctx context.Context, payload SMAPayload) (*timeserie.TimeSerie[float64], error)
	ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	Close(ctx context.Context)
}

type SMAPayload struct {
	Exchange     string
	Pair         string
	Period       period.Symbol
	Start        time.Time
	End          time.Time
	PeriodNumber uint
	PriceType    candlestick.PriceType
}
