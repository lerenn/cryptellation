package retry

import (
	"context"

	common "github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/models/account"

	client "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

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

func (r retry) CreateForwardTest(ctx context.Context, payload forwardtest.NewPayload) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		id, err = r.client.CreateForwardTest(ctx, payload)
		return err
	})
	return id, err
}

func (r retry) ListForwardTests(ctx context.Context) ([]forwardtest.ForwardTest, error) {
	var fts []forwardtest.ForwardTest
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		fts, err = r.client.ListForwardTests(ctx)
		return err
	})
	return fts, err
}

func (r retry) CreateOrder(ctx context.Context, payload common.OrderCreationPayload) error {
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		err = r.client.CreateOrder(ctx, payload)
		return err
	})
	return err
}

func (r retry) GetAccounts(ctx context.Context, forwardTestID uuid.UUID) (map[string]account.Account, error) {
	var accounts map[string]account.Account
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		accounts, err = r.client.GetAccounts(ctx, forwardTestID)
		return err
	})
	return accounts, err
}

func (r retry) GetStatus(ctx context.Context, forwardTestID uuid.UUID) (forwardtest.Status, error) {
	var status forwardtest.Status
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		status, err = r.client.GetStatus(ctx, forwardTestID)
		return err
	})
	return status, err
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
