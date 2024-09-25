package client

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/client"
	"github.com/lerenn/cryptellation/pkg/config"

	backtests "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	backtestsnats "github.com/lerenn/cryptellation/svc/backtests/clients/go/nats"

	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	candlesticksnats "github.com/lerenn/cryptellation/svc/candlesticks/clients/go/nats"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"

	exchanges "github.com/lerenn/cryptellation/svc/exchanges/clients/go"
	exchangesnats "github.com/lerenn/cryptellation/svc/exchanges/clients/go/nats"

	forwardtests "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
	forwardtestsnats "github.com/lerenn/cryptellation/svc/forwardtests/clients/go/nats"

	indicators "github.com/lerenn/cryptellation/svc/indicators/clients/go"
	indicatorsnats "github.com/lerenn/cryptellation/svc/indicators/clients/go/nats"

	ticks "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	ticksnats "github.com/lerenn/cryptellation/svc/ticks/clients/go/nats"
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
	svc.Backtests, err = backtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	// Set candlesticks
	svc.Candlesticks, err = candlesticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	// Set exchanges
	svc.Exchanges, err = exchangesnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	// Set forward tests
	svc.Forwardtests, err = forwardtestsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	// Set indicators
	svc.Indicators, err = indicatorsnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	// Set ticks
	svc.Ticks, err = ticksnats.New(c)
	if err != nil {
		svc.Close(context.TODO())
		return
	}

	return
}

func (c Services) ExportBacktestsData(ctx context.Context, backtestID uuid.UUID) (RunDataExport, error) {
	// Get backtest info
	bt, err := c.Backtests.Get(ctx, backtestID)
	if err != nil {
		return RunDataExport{}, err
	}

	// Get candlesticks
	csList := make(map[string]map[string][]candlestick.Candlestick)
	for _, ts := range bt.TickSubscriptions {
		cs, err := c.Candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
			Exchange: ts.Exchange,
			Pair:     ts.Pair,
			Period:   bt.PeriodBetweenEvents,
			Start:    &bt.StartTime,
			End:      &bt.EndTime,
		})
		if err != nil {
			return RunDataExport{}, err
		}

		// Check if the map for exchange is initialized
		if _, ok := csList[ts.Exchange]; !ok {
			csList[ts.Exchange] = make(map[string][]candlestick.Candlestick)
		}

		// Check if the map for pair is initialized
		if _, ok := csList[ts.Exchange][ts.Pair]; !ok {
			csList[ts.Exchange][ts.Pair] = make([]candlestick.Candlestick, 0)
		}

		csList[ts.Exchange][ts.Pair] = cs.ToArray()
	}

	// Get orders
	orders, err := c.Backtests.ListOrders(ctx, backtestID)
	if err != nil {
		return RunDataExport{}, err
	}

	return RunDataExport{
		ID:           backtestID,
		Mode:         ModeIsBacktest,
		Candlesticks: csList,
		Orders:       orders,
	}, nil
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
