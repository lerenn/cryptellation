package candlesticks

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers"
)

type ServiceConfig struct {
	Adapters    adapters.Config
	Controllers controllers.Config
}

type Service struct {
	controllers []controllers.Controller
}

func New(c ServiceConfig) (Service, error) {
	// Set adapters
	a, err := adapters.New(c.Adapters)
	if err != nil {
		return Service{}, err
	}

	// Set components
	candlesticks := candlesticks.New(a.Database, a.Exchanges)

	// Set controllers
	controllers, err := controllers.New(c.Controllers, candlesticks)
	if err != nil {
		return Service{}, err
	}

	return Service{
		controllers: controllers,
	}, nil
}

func (s Service) Serve() {
	for _, c := range s.controllers {
		c.Listen()
	}
}

func (s Service) Close() {
	for _, c := range s.controllers {
		c.Close()
	}
}
