package client

import (
	"context"
	"errors"
	"sync"

	"github.com/lerenn/cryptellation/pkg/client"
)

func (c Client) ServicesInfo(ctx context.Context) (servicesInfo map[string]client.ServiceInfo, err error) {
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

	wg.Add(5)
	go cb("backtests", c.backtests)
	go cb("candlesticks", c.candlesticks)
	go cb("exchanges", c.exchanges)
	go cb("indicators", c.indicators)
	go cb("ticks", c.ticks)

	wg.Wait()

	return
}
