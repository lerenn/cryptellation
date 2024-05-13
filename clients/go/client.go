package client

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/config"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	natsbacktests "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	natscandlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	natsexchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"
	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	natsindicators "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"
	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	natsticks "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
)

type Client struct {
	backtests    backtests.Client
	candlesticks candlesticks.Client
	exchanges    exchanges.Client
	indicators   indicators.Client
	ticks        ticks.Client
}

func NewNATSClient(c config.NATS) (client Client, err error) {
	client.backtests, err = natsbacktests.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}

	client.candlesticks, err = natscandlesticks.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}

	client.exchanges, err = natsexchanges.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}

	client.indicators, err = natsindicators.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}

	client.ticks, err = natsticks.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}

	return
}

func (c Client) Backtests() backtests.Client {
	return c.backtests
}

func (c Client) Candlesticks() candlesticks.Client {
	return c.candlesticks
}

func (c Client) Exchanges() exchanges.Client {
	return c.exchanges
}

func (c Client) Indicators() indicators.Client {
	return c.indicators
}

func (c Client) Ticks() ticks.Client {
	return c.ticks
}

func (c *Client) Close(ctx context.Context) {
	if c.backtests != nil {
		c.backtests.Close(ctx)
	}

	if c.candlesticks != nil {
		c.candlesticks.Close(ctx)
	}

	if c.exchanges != nil {
		c.exchanges.Close(ctx)
	}

	if c.indicators != nil {
		c.indicators.Close(ctx)
	}

	if c.ticks != nil {
		c.ticks.Close(ctx)
	}
}
