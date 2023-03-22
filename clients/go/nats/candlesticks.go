package nats

import (
	"context"
	"fmt"
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/candlesticks/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/types/candlestick"
	"github.com/digital-feather/cryptellation/pkg/utils"
	"github.com/nats-io/nats.go"
)

type Candlesticks struct {
	nats *nats.Conn
	ctrl *generated.ClientController
}

func NewCandlesticks(c config.NATS) (client.Candlesticks, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := generated.NewClientController(generated.NewNATSController(conn))
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
	reqMsg := generated.NewCandlesticksListRequestMessage()
	reqMsg.Payload.ExchangeName = generated.ExchangeNameSchema(payload.ExchangeName)
	reqMsg.Payload.PairSymbol = generated.PairSymbolSchema(payload.PairSymbol)
	reqMsg.Payload.PeriodSymbol = generated.PeriodSymbolSchema(payload.Period.String())
	if payload.Start != nil {
		reqMsg.Payload.Start = utils.ToReference(generated.DateSchema(*payload.Start))
	}
	if payload.End != nil {
		reqMsg.Payload.End = utils.ToReference(generated.DateSchema(*payload.End))
	}
	reqMsg.Payload.Limit = generated.LimitSchema(payload.Limit)

	// Send request
	respMsg, err := c.ctrl.WaitForCandlesticksListResponse(ctx, reqMsg, func() error {
		return c.ctrl.PublishCandlesticksListRequest(reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To candlestick list
	list := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	})
	for _, c := range *respMsg.Payload.Candlesticks {
		if err := list.Set(time.Time(c.Time), candlestick.Candlestick{
			Open:   c.Open,
			High:   c.High,
			Low:    c.Low,
			Close:  c.Close,
			Volume: c.Volume,
		}); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (c Candlesticks) Close() {
	c.ctrl.Close()
	c.nats.Close()
}
