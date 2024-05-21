package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
)

var (
	ErrUnsupportedLiveMode    = errors.New("unsupported live mode")
	ErrUnsupportedForwardMode = errors.New("unsupported forward mode")
)

type Run struct {
	ID       uuid.UUID
	Mode     Mode
	Services Services
	Time     time.Time
}

func (r Run) CreateOrder(ctx context.Context, payload client.OrderCreationPayload) error {
	payload.BacktestID = r.ID

	switch r.Mode {
	case ModeIsBacktest:
		return r.Services.Backtests().CreateOrder(ctx, payload)
	case ModeIsLive:
		return ErrUnsupportedLiveMode
	case ModeIsForward:
		return ErrUnsupportedForwardMode
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedMode, r.Mode)
	}
}

func (r Run) GetAccounts(ctx context.Context) (map[string]account.Account, error) {
	switch r.Mode {
	case ModeIsBacktest:
		return r.Services.Backtests().GetAccounts(ctx, r.ID)
	case ModeIsLive:
		return nil, ErrUnsupportedLiveMode
	case ModeIsForward:
		return nil, ErrUnsupportedForwardMode
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedMode, r.Mode)
	}
}
