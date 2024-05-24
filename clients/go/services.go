package client

import (
	"context"
	"errors"
	"sync"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"
	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	natsbacktests "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	natscandlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	natsexchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"
	forwardtests "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	natsforwardtests "github.com/lerenn/cryptellation/svc/forwardtests/clients/go/nats"
	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	natsindicators "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"
	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	natsticks "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
)

type Services struct {
	backtests    backtests.Client
	candlesticks candlesticks.Client
	exchanges    exchanges.Client
	forwardtests forwardtests.Client
	indicators   indicators.Client
	ticks        ticks.Client
}

func NewServices(c config.NATS) (client Services, err error) {
	backtests, err := natsbacktests.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.backtests = backtests

	candlesticks, err := natscandlesticks.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.candlesticks = candlesticks

	exchanges, err := natsexchanges.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.exchanges = exchanges

	forwardtests, err := natsforwardtests.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.forwardtests = forwardtests

	indicators, err := natsindicators.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.indicators = indicators

	ticks, err := natsticks.NewClient(c)
	if err != nil {
		client.Close(context.TODO())
		return
	}
	client.ticks = ticks

	return
}

func (c Services) Backtests() backtests.Client {
	return c.backtests
}

func (c Services) Candlesticks() candlesticks.Client {
	return c.candlesticks
}

func (c Services) Exchanges() exchanges.Client {
	return c.exchanges
}

func (c Services) ForwardTests() forwardtests.Client {
	return c.forwardtests
}

func (c Services) Indicators() indicators.Client {
	return c.indicators
}

func (c Services) Ticks() ticks.Client {
	return c.ticks
}

func (c Services) ServicesInfo(ctx context.Context) (servicesInfo map[string]client.ServiceInfo, err error) {
	servicesInfo = make(map[string]client.ServiceInfo)
	lock := sync.Mutex{}

	type ServiceInfoer interface {
		ServiceInfo(ctx context.Context) (client.ServiceInfo, error)
	}

	var wg sync.WaitGroup
	cb := func(service string, cl ServiceInfoer) {
		var localErr error

		defer wg.Done()
		lock.Lock()
		defer lock.Unlock()

		servicesInfo[service], localErr = cl.ServiceInfo(ctx)
		if err != nil {
			err = errors.Join(err, localErr)
		}
	}

	wg.Add(6)
	go cb("backtests", c.backtests)
	go cb("candlesticks", c.candlesticks)
	go cb("exchanges", c.exchanges)
	go cb("forwardtests", c.forwardtests)
	go cb("indicators", c.indicators)
	go cb("ticks", c.ticks)

	wg.Wait()

	return
}

func (c *Services) Close(ctx context.Context) {
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
