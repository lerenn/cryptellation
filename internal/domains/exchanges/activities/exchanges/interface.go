// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=interface.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
	"go.temporal.io/sdk/worker"
)

// GetExchangeInfoActivityName is the name of the GetExchangeInfo activity.
const GetExchangeInfoActivityName = "GetExchangeInfoActivity"

type (
	// GetExchangeInfoParams is the parameters for the GetExchangeInfo activity.
	GetExchangeInfoParams struct {
		Name string
	}

	// GetExchangeInfoResult is the result for the GetExchangeInfo activity.
	GetExchangeInfoResult struct {
		Exchange exchange.Exchange
	}
)

// Interface is the interface for the exchanges activities.
type Interface interface {
	Register(w worker.Worker)

	GetExchangeInfo(ctx context.Context, params GetExchangeInfoParams) (GetExchangeInfoResult, error)
}
