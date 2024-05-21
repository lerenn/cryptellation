package domain

import (
	"context"

	"github.com/google/uuid"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

// Test interface implementation
var _ app.ForwardTests = (*ForwardTests)(nil)

type ForwardTests struct {
	db           db.Port
	candlesticks candlesticks.Client
}

func New(db db.Port, csClient candlesticks.Client) *ForwardTests {
	if db == nil {
		panic("nil db")
	}

	if csClient == nil {
		panic("nil candlesticks client")
	}

	return &ForwardTests{
		db:           db,
		candlesticks: csClient,
	}
}

func (f ForwardTests) Create(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error) {
	if err := payload.Validate(); err != nil {
		return uuid.Nil, err
	}

	ft := forwardtest.New(payload)
	return ft.ID, f.db.CreateForwardTest(ctx, ft)
}
