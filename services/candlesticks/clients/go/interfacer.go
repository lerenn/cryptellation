// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=interfacer.go -destination=mock.gen.go -package client

package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/candlesticks/clients/go/payloads"
	"github.com/digital-feather/cryptellation/services/candlesticks/pkg/models/candlestick"
)

type Interfacer interface {
	ReadCandlesticks(ctx context.Context, payload payloads.ReadCandlesticksPayload) (*candlestick.List, error)
}
