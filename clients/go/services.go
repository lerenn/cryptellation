package client

import (
	"context"
	"errors"
	"sync"

	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"

	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	backtestsnats "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"
	backtestsretry "github.com/lerenn/cryptellation/svc/backtests/clients/go/retry"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlestickscache "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/cache"
	candlesticksnats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	candlesticksretry "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/retry"

	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	exchangescache "github.com/lerenn/cryptellation/svc/exchanges/clients/go/cache"
	exchangesnats "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"
	exchangesretry "github.com/lerenn/cryptellation/svc/exchanges/clients/go/retry"

	forwardtests "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	forwardtestsnats "github.com/lerenn/cryptellation/svc/forwardtests/clients/go/nats"
	forwardtestsretry "github.com/lerenn/cryptellation/svc/forwardtests/clients/go/retry"

	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	indicatorscache "github.com/lerenn/cryptellation/svc/indicators/clients/go/cache"
	indicatorsnats "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"
	indicatorsretry "github.com/lerenn/cryptellation/svc/indicators/clients/go/retry"

	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	ticksnats "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
	ticksretry "github.com/lerenn/cryptellation/svc/ticks/clients/go/retry"
)

type Services struct {
	Backtests    backtests.Client
	Candlesticks candlesticks.Client
	Exchanges    exchanges.Client
	Forwardtests forwardtests.Client
	Indicators   indicators.Client
	Ticks        ticks.Client
}

func NewServices(c config.NATS) (svc Services, err error) {
	// Set backtests
	backtests, err := backtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	backtests = backtestsretry.New(backtests)
	svc.Backtests = backtests

	// Set candlesticks
	candlesticks, err := candlesticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	candlesticks = candlesticksretry.New(candlesticks)
	candlesticks = candlestickscache.New(candlesticks)
	svc.Candlesticks = candlesticks

	// Set exchanges
	exchanges, err := exchangesnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	exchanges = exchangesretry.New(exchanges)
	exchanges = exchangescache.New(exchanges)
	svc.Exchanges = exchanges

	// Set forward tests
	forwardtests, err := forwardtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	forwardtests = forwardtestsretry.New(forwardtests)
	svc.Forwardtests = forwardtests

	// Set indicators
	indicators, err := indicatorsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	indicators = indicatorscache.New(indicators)
	indicators = indicatorsretry.New(indicators)
	svc.Indicators = indicators

	// Set ticks
	ticks, err := ticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}
	ticks = ticksretry.New(ticks)
	svc.Ticks = ticks

	return
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
	go cb("backtests", c.Backtests)
	go cb("candlesticks", c.Candlesticks)
	go cb("exchanges", c.Exchanges)
	go cb("forwardtests", c.Forwardtests)
	go cb("indicators", c.Indicators)
	go cb("ticks", c.Ticks)

	wg.Wait()

	return
}

func (c *Services) Close(ctx context.Context) {
	if c.Backtests != nil {
		c.Backtests.Close(ctx)
	}

	if c.Candlesticks != nil {
		c.Candlesticks.Close(ctx)
	}

	if c.Exchanges != nil {
		c.Exchanges.Close(ctx)
	}

	if c.Indicators != nil {
		c.Indicators.Close(ctx)
	}

	if c.Ticks != nil {
		c.Ticks.Close(ctx)
	}
}
