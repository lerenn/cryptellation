// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=interface.go -destination=mock.gen.go -package db

package db

import (
	"context"

	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
)

// CreateExchangesActivityName is the name of the CreateExchanges activity.
const CreateExchangesActivityName = "CreateExchangesActivity"

type (
	// CreateExchangesParams is the parameters for the CreateExchanges activity.
	CreateExchangesParams struct {
		Exchanges []exchange.Exchange
	}

	// CreateExchangesResult is the result for the CreateExchanges activity.
	CreateExchangesResult struct{}
)

// ReadExchangesActivityName is the name of the ReadExchanges activity.
const ReadExchangesActivityName = "ReadExchangesActivity"

type (
	// ReadExchangesParams is the parameters for the ReadExchanges activity.
	ReadExchangesParams struct {
		Names []string
	}

	// ReadExchangesResult is the result for the ReadExchanges activity.
	ReadExchangesResult struct {
		Exchanges []exchange.Exchange
	}
)

// UpdateExchangesActivityName is the name of the UpdateExchanges activity.
const UpdateExchangesActivityName = "UpdateExchangesActivity"

type (
	// UpdateExchangesParams is the parameters for the UpdateExchanges activity.
	UpdateExchangesParams struct {
		Exchanges []exchange.Exchange
	}

	// UpdateExchangesResult is the result for the UpdateExchanges activity.
	UpdateExchangesResult struct{}
)

// DeleteExchangesActivityName is the name of the DeleteExchanges activity.
const DeleteExchangesActivityName = "DeleteExchangesActivity"

type (
	// DeleteExchangesParams is the parameters for the DeleteExchanges activity.
	DeleteExchangesParams struct {
		Names []string
	}

	// DeleteExchangesResult is the result for the DeleteExchanges activity.
	DeleteExchangesResult struct{}
)

// Interface is the interface that the database activities must implement
type Interface interface {
	CreateExchanges(ctx context.Context, params CreateExchangesParams) (CreateExchangesResult, error)
	ReadExchanges(ctx context.Context, params ReadExchangesParams) (ReadExchangesResult, error)
	UpdateExchanges(ctx context.Context, params UpdateExchangesParams) (UpdateExchangesResult, error)
	DeleteExchanges(ctx context.Context, params DeleteExchangesParams) (DeleteExchangesResult, error)
}
