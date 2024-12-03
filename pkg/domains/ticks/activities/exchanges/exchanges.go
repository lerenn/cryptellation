// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"go.temporal.io/sdk/worker"
)

// ListenSymbolActivityName is the name of the activity to listen to a symbol.
const ListenSymbolActivityName = "ListenSymbolActivity"

type (
	// ListenSymbolParams is the parameters for the ListenSymbolActivity.
	ListenSymbolParams struct {
		ParentWorkflowID string
		Exchange         string
		Symbol           string
	}

	// ListenSymbolResults is the results for the ListenSymbolActivity.
	ListenSymbolResults struct{}
)

// Exchanges is the exchanges activities for ticks.
type Exchanges interface {
	Register(w worker.Worker)

	ListenSymbolActivity(ctx context.Context, params ListenSymbolParams) (ListenSymbolResults, error)
}
