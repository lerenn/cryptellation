package nats

import (
	"context"
	"fmt"

	"github.com/lerenn/asyncapi-codegen/pkg/log"
	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/internal/ctrl/indicators/events"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/timeserie"
	"github.com/nats-io/nats.go"
)

type Indicators struct {
	nats *nats.Conn
	ctrl *events.ClientController
}

func NewIndicators(c config.NATS) (client.Indicators, error) {
	conn, err := nats.Connect(c.URL())
	if err != nil {
		return nil, err
	}

	ctrl, err := events.NewClientController(events.NewNATSController(conn))
	if err != nil {
		return nil, err
	}
	ctrl.SetLogger(log.NewECS())

	return Indicators{
		nats: conn,
		ctrl: ctrl,
	}, nil
}

func (ex Indicators) SMA(ctx context.Context, payload client.SMAPayload) (*timeserie.TimeSerie[float64], error) {
	// Set message
	reqMsg := events.NewSmaRequestMessage()
	reqMsg.Set(payload)

	// Send request
	respMsg, err := ex.ctrl.WaitForCryptellationIndicatorsSmaResponse(ctx, reqMsg, func(ctx context.Context) error {
		return ex.ctrl.PublishCryptellationIndicatorsSmaRequest(ctx, reqMsg)
	})
	if err != nil {
		return nil, err
	}

	// Check error
	if respMsg.Payload.Error != nil {
		return nil, fmt.Errorf("%d Error: %s", respMsg.Payload.Error.Code, respMsg.Payload.Error.Message)
	}

	// To indicator list
	return respMsg.ToModel(), nil
}

func (ex Indicators) Close(ctx context.Context) {
	ex.ctrl.Close(ctx)
	ex.nats.Close()
}
