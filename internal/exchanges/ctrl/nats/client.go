package nats

import (
	"context"
	"fmt"

	"github.com/digital-feather/cryptellation/internal/exchanges/ctrl/nats/internal"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/exchange"
	"github.com/nats-io/nats.go"
)

type client struct {
	nats *nats.Conn
	ctrl *internal.ClientController
}

func New(c config.NATS) (Client, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := internal.NewClientController(internal.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return client{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (c client) ReadExchanges(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := internal.NewExchangesRequestMessage()
	reqMsg.Payload = make([]internal.ExchangeNameSchema, 0, len(names))
	for _, name := range names {
		reqMsg.Payload = append(reqMsg.Payload, internal.ExchangeNameSchema(name))
	}

	// Send request
	respMsg, err := c.ctrl.WaitForExchangesListResponse(ctx, reqMsg, func() error {
		return c.ctrl.PublishExchangesListRequest(reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To exchange list
	exchanges := make([]exchange.Exchange, len(respMsg.Payload.Exchanges))
	for i, exch := range respMsg.Payload.Exchanges {
		// Periods
		periods := make([]string, len(exch.Periods))
		for j, p := range exch.Periods {
			periods[j] = string(p)
		}

		// Pairs
		pairs := make([]string, len(exch.Pairs))
		for j, p := range exch.Pairs {
			pairs[j] = string(p)
		}

		exchanges[i] = exchange.Exchange{
			Name:           string(exch.Name),
			Fees:           exch.Fees,
			PairsSymbols:   pairs,
			PeriodsSymbols: periods,
			LastSyncTime:   exch.LastSyncTime,
		}
	}

	return exchanges, nil
}

func (c client) Close() {
	c.ctrl.Close()
	c.nats.Close()
}
