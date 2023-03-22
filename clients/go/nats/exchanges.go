package nats

import (
	"context"
	"fmt"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/exchanges/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/exchange"
	"github.com/nats-io/nats.go"
)

type Exchanges struct {
	nats *nats.Conn
	ctrl *generated.ClientController
}

func NewExchanges(c config.NATS) (client.Exchanges, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := generated.NewClientController(generated.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return Exchanges{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (ex Exchanges) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	// Set message
	reqMsg := generated.NewExchangesRequestMessage()
	reqMsg.Payload = make([]generated.ExchangeNameSchema, 0, len(names))
	for _, name := range names {
		reqMsg.Payload = append(reqMsg.Payload, generated.ExchangeNameSchema(name))
	}

	// Send request
	respMsg, err := ex.ctrl.WaitForExchangesListResponse(ctx, reqMsg, func() error {
		return ex.ctrl.PublishExchangesListRequest(reqMsg)
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

func (ex Exchanges) Close() {
	ex.ctrl.Close()
	ex.nats.Close()
}
