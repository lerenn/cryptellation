package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"

	"github.com/google/uuid"
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
		return r.Services.BacktestsClient().CreateOrder(ctx, payload)
	case ModeIsForwardTest:
		return r.Services.ForwardTestsClient().CreateOrder(ctx, payload)
	default:
		return fmt.Errorf("%w for CreateOrder(): %q", ErrUnsupportedMode, r.Mode)
	}
}

func (r Run) GetAccounts(ctx context.Context) (map[string]account.Account, error) {
	switch r.Mode {
	case ModeIsBacktest:
		return r.Services.BacktestsClient().GetAccounts(ctx, r.ID)
	case ModeIsForwardTest:
		return r.Services.ForwardTestsClient().GetAccounts(ctx, r.ID)
	default:
		return nil, fmt.Errorf("%w for GetAccounts(): %q", ErrUnsupportedMode, r.Mode)
	}
}
