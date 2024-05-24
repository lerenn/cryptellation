package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
)

type Run struct {
	ID       uuid.UUID
	Mode     Mode
	Services Services
	Time     time.Time
}

func (r Run) CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error {
	payload.RunID = r.ID

	switch r.Mode {
	case ModeIsBacktest:
		return r.Services.Backtests().CreateOrder(ctx, payload)
	case ModeIsForwardTest:
		return r.Services.ForwardTests().CreateOrder(ctx, payload)
	default:
		return fmt.Errorf("%w for CreateOrder(): %q", ErrUnsupportedMode, r.Mode)
	}
}

func (r Run) GetAccounts(ctx context.Context) (map[string]account.Account, error) {
	switch r.Mode {
	case ModeIsBacktest:
		return r.Services.Backtests().GetAccounts(ctx, r.ID)
	case ModeIsForwardTest:
		return r.Services.ForwardTests().GetAccounts(ctx, r.ID)
	default:
		return nil, fmt.Errorf("%w for GetAccounts(): %q", ErrUnsupportedMode, r.Mode)
	}
}
