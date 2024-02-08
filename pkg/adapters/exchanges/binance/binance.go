package binance

import (
	"io"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/svc/exchanges/pkg/exchange"
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

func New() (*Service, error) {
	// Get config
	config := config.LoadBinance()

	c := client.NewClient(config.ApiKey, config.SecretKey)
	c.Logger.SetOutput(io.Discard)

	// Return service
	return &Service{
		Client: c,
	}, config.Validate()
}
