package nats

import (
	"context"
	"time"

	"github.com/digital-feather/cryptellation/internal/candlesticks/ctrl/nats/internal"
	"github.com/digital-feather/cryptellation/pkg/candlestick"
	"github.com/digital-feather/cryptellation/pkg/config"
	"github.com/digital-feather/cryptellation/pkg/utils"
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

func (c client) ReadCandlesticks(ctx context.Context, payload ReadCandlesticksPayload) (*candlestick.List, error) {
	// Set message
	reqMsg := internal.NewCandlesticksListRequestMessage()
	reqMsg.Payload.ExchangeName = internal.ExchangeNameSchema(payload.ExchangeName)
	reqMsg.Payload.PairSymbol = internal.PairSymbolSchema(payload.PairSymbol)
	reqMsg.Payload.PeriodSymbol = internal.PeriodSymbolSchema(payload.Period.String())
	if payload.Start != nil {
		reqMsg.Payload.Start = utils.ToReference(internal.DateSchema(*payload.Start))
	}
	if payload.End != nil {
		reqMsg.Payload.End = utils.ToReference(internal.DateSchema(*payload.End))
	}
	reqMsg.Payload.Limit = internal.LimitSchema(payload.Limit)

	// Send request
	respMsg, err := c.ctrl.WaitForCandlesticksListResponse(reqMsg, func() error {
		return c.ctrl.PublishCandlesticksListRequest(reqMsg)
	}, time.Second)
	if err != nil {
		return nil, err
	}

	// To candlestick list
	list := candlestick.NewEmptyList(candlestick.ListID{
		ExchangeName: payload.ExchangeName,
		PairSymbol:   payload.PairSymbol,
		Period:       payload.Period,
	})
	for _, c := range *respMsg.Payload.Candlesticks {
		list.Set(time.Time(c.Time), candlestick.Candlestick{
			Open:   c.Open,
			High:   c.High,
			Low:    c.Low,
			Close:  c.Close,
			Volume: c.Volume,
		})
	}

	return nil, nil
}

func (c client) Close() {
	c.ctrl.Close()
	c.nats.Close()
}
