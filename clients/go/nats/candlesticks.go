package nats

import (
	"context"
	"fmt"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/ctrl/candlesticks/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
	"github.com/nats-io/nats.go"
)

type Candlesticks struct {
	nats *nats.Conn
	ctrl *events.ClientController
}

func NewCandlesticks(c config.NATS) (client.Candlesticks, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := events.NewClientController(events.NewNATSController(conn))
	if err != nil {
		return nil, err
	}

	return Candlesticks{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (c Candlesticks) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	// Set message
	reqMsg := events.NewCandlesticksListRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := c.ctrl.WaitForCryptellationCandlesticksListResponse(ctx, reqMsg, func() error {
		return c.ctrl.PublishCryptellationCandlesticksListRequest(reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To candlestick list
	return respMsg.ToModel(payload.ExchangeName, payload.PairSymbol, payload.Period)
}

func (c Candlesticks) Close() {
	c.ctrl.Close()
	c.nats.Close()
}
