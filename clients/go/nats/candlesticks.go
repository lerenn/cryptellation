package nats

import (
	"context"
	"fmt"

	asyncapi "github.com/digital-feather/cryptellation/api/asyncapi/candlesticks"
	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/nats-io/nats.go"
)

type Candlesticks struct {
	nats *nats.Conn
	ctrl *asyncapi.ClientController
}

func NewCandlesticks(c config.NATS) (client.Candlesticks, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := asyncapi.NewClientController(asyncapi.NewNATSController(conn))
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
	reqMsg := asyncapi.NewCandlesticksListRequestMessage()
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
