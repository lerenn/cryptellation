package indicators

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/models/timeserie"
)

type Interface interface {
	GetCachedSMA(ctx context.Context, payload GetCachedSMAPayload) (*timeserie.TimeSerie[float64], error)
}
