package nats

import (
	"context"
	"fmt"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/ctrl/exchanges/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/exchange"
	"github.com/nats-io/nats.go"
)

type Exchanges struct {
	nats *nats.Conn
	ctrl *events.ClientController
}

func NewExchanges(c config.NATS) (client.Exchanges, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := events.NewClientController(events.NewNATSController(conn))
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
	reqMsg := events.NewExchangesRequestMessage()
	reqMsg.Set(names...)

	// Send request
	respMsg, err := ex.ctrl.WaitForCryptellationExchangesListResponse(ctx, reqMsg, func() error {
		return ex.ctrl.PublishCryptellationExchangesListRequest(reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To exchange list
	return respMsg.ToModel(), nil
}

func (ex Exchanges) Close() {
	ex.ctrl.Close()
	ex.nats.Close()
}
