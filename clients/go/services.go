package client

import (
	"context"
	"errors"
	"sync"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"

	backtests "github.com/lerenn/cryptellation/client/clients/go"
	backtestsnats "github.com/lerenn/cryptellation/client/clients/go/nats"
	backtestsretry "github.com/lerenn/cryptellation/client/clients/go/retry"

	candlesticks "github.com/lerenn/cryptellation/candlesticks/clients/go"
	candlestickscache "github.com/lerenn/cryptellation/candlesticks/clients/go/cache"
	candlesticksnats "github.com/lerenn/cryptellation/candlesticks/clients/go/nats"
	candlesticksretry "github.com/lerenn/cryptellation/candlesticks/clients/go/retry"

	exchanges "github.com/lerenn/cryptellation/exchanges/clients/go"
	exchangescache "github.com/lerenn/cryptellation/exchanges/clients/go/cache"
	exchangesnats "github.com/lerenn/cryptellation/exchanges/clients/go/nats"
	exchangesretry "github.com/lerenn/cryptellation/exchanges/clients/go/retry"

	forwardtests "github.com/lerenn/cryptellation/forwardtests/clients/go"
	forwardtestsnats "github.com/lerenn/cryptellation/forwardtests/clients/go/nats"
	forwardtestsretry "github.com/lerenn/cryptellation/forwardtests/clients/go/retry"

	indicators "github.com/lerenn/cryptellation/indicators/clients/go"
	indicatorscache "github.com/lerenn/cryptellation/indicators/clients/go/cache"
	indicatorsnats "github.com/lerenn/cryptellation/indicators/clients/go/nats"
	indicatorsretry "github.com/lerenn/cryptellation/indicators/clients/go/retry"

	ticks "github.com/lerenn/cryptellation/ticks/clients/go"
	ticksnats "github.com/lerenn/cryptellation/ticks/clients/go/nats"
	ticksretry "github.com/lerenn/cryptellation/ticks/clients/go/retry"
)

type Services struct {
	backtests    backtests.Client
	candlesticks candlesticks.Client
	exchanges    exchanges.Client
	forwardtests forwardtests.Client
	indicators   indicators.Client
	ticks        ticks.Client
}

func NewServices(c config.NATS) (svc Services, err error) {
	// Set backtests
	backtests, err := backtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	backtests = backtestsretry.New(backtests)
	svc.backtests = backtests

	// Set candlesticks
	candlesticks, err := candlesticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	candlesticks = candlesticksretry.New(candlesticks)
	candlesticks = candlestickscache.New(candlesticks)
	svc.candlesticks = candlesticks

	// Set exchanges
	exchanges, err := exchangesnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	exchanges = exchangesretry.New(exchanges)
	exchanges = exchangescache.New(exchanges)
	svc.exchanges = exchanges

	// Set forward tests
	forwardtests, err := forwardtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	forwardtests = forwardtestsretry.New(forwardtests)
	svc.forwardtests = forwardtests

	// Set indicators
	indicators, err := indicatorsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	indicators = indicatorscache.New(indicators)
	indicators = indicatorsretry.New(indicators)
	svc.indicators = indicators

	// Set ticks
	ticks, err := ticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	ticks = ticksretry.New(ticks)
	svc.ticks = ticks

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
