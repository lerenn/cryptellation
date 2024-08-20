package retry

import (
	"context"

	common "github.com/lerenn/cryptellation/pkg/client"

	client "github.com/lerenn/cryptellation/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/candlesticks/pkg/candlestick"
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

func (r retry) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (resp *candlestick.List, err error) {
	err = r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		resp, err = r.client.Read(ctx, payload)
		return err
	})
	return
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
