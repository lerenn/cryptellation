// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=exchanges.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"go.temporal.io/sdk/worker"
)

const ListenSymbolActivityName = "ListenSymbolActivity"

type (
	ListenSymbolParams struct {
		ParentWorkflowID string
		Exchange         string
		Symbol           string
	}

	ListenSymbolResults struct{}
)

type Exchanges interface {
	Register(w worker.Worker)

	ListenSymbolActivity(ctx context.Context, params ListenSymbolParams) (ListenSymbolResults, error)
}
