package binance

import (
	"io"

	"cryptellation/internal/config"

	"cryptellation/svc/exchanges/pkg/exchange"

	client "github.com/adshao/go-binance/v2"
)

var (
	Infos = exchange.Exchange{
		Name: "binance",
		Periods: []string{
			"M1", "M3", "M5", "M15", "M30",
			"H1", "H2", "H4", "H6", "H8", "H12",
			"D1", "D3",
			"W1",
		},
		Fees: 0.1,
	}
)

type Service struct {
	Client *client.Client
}

func New(conf config.Binance) (*Service, error) {
	c := client.NewClient(conf.ApiKey, conf.SecretKey)
	c.Logger.SetOutput(io.Discard)

	// Return service
	return &Service{
		Client: c,
	}, conf.Validate()
}
