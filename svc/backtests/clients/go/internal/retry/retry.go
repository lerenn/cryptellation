package retry

import (
	"context"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/order"

	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/google/uuid"
)

type retry struct {
	client client.Client
	common.Retryable
}

func New(client client.Client, options ...option) client.Client {
	r := retry{
		client:    client,
		Retryable: common.DefaultRetryable,
	}

	// Execute options
	for _, option := range options {
		option(&r)
	}

	return &r
}

func (r retry) Advance(ctx context.Context, backtestID uuid.UUID) error {
	return r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		return r.client.Advance(ctx, backtestID)
	})
}

func (r retry) Create(ctx context.Context, payload client.BacktestCreationPayload) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		id, err = r.client.Create(ctx, payload)
		return err
	})
	return id, err
}

func (r retry) CreateOrder(ctx context.Context, payload common.OrderCreationPayload) error {
	return r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		return r.client.CreateOrder(ctx, payload)
	})

}

func (r retry) Get(ctx context.Context, backtestID uuid.UUID) (backtest.Backtest, error) {
	var bt backtest.Backtest
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		bt, err = r.client.Get(ctx, backtestID)
		return err
	})
	return bt, err
}

func (r retry) GetAccounts(ctx context.Context, backtestID uuid.UUID) (map[string]account.Account, error) {
	var accounts map[string]account.Account
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		accounts, err = r.client.GetAccounts(ctx, backtestID)
		return err
	})
	return accounts, err

}

func (r retry) Subscribe(ctx context.Context, backtestID uuid.UUID, exchange, pair string) error {
	return r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		return r.client.Subscribe(ctx, backtestID, exchange, pair)
	})
}

func (r retry) ListenEvents(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error) {
	var events <-chan event.Event
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		events, err = r.client.ListenEvents(ctx, backtestID)
		return err
	})
	return events, err
}

func (r retry) List(ctx context.Context) ([]backtest.Backtest, error) {
	var backtests []backtest.Backtest
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		backtests, err = r.client.List(ctx)
		return err
	})
	return backtests, err
}

func (r retry) ListOrders(ctx context.Context, backtestID uuid.UUID) ([]order.Order, error) {
	var orders []order.Order
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		orders, err = r.client.ListOrders(ctx, backtestID)
		return err
	})
	return orders, err
}

func (r retry) ServiceInfo(ctx context.Context) (resp common.ServiceInfo, err error) {
	err = r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		resp, err = r.client.ServiceInfo(ctx)
		return err
	})
	return
}

func (r retry) Close(ctx context.Context) {
	r.client.Close(ctx)
}
