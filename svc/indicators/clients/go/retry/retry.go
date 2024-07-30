package retry

import (
	"context"

	common "cryptellation/pkg/client"
	"cryptellation/pkg/models/timeserie"

	client "cryptellation/svc/indicators/clients/go"
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

func (r retry) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	var serie *timeserie.TimeSerie[float64]
	err := r.Retryable.Exec(ctx, func(ctx context.Context) (err error) {
		serie, err = r.client.SMA(ctx, payload)
		return err
	})
	return serie, err

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
