package activities

import (
	"io"

	client "github.com/adshao/go-binance/v2"
	"github.com/lerenn/cryptellation/v1/pkg/config"
	"github.com/lerenn/cryptellation/v1/pkg/models/exchange"
)

var (
	// BinanceInfos represents the Binance exchange informations.
	BinanceInfos = exchange.Exchange{
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

// Binance represents a Binance accessor.
type Binance struct {
	Client *client.Client
}

// NewBinance creates a new Binance accessor.
func NewBinance(conf config.Binance) (*Binance, error) {
	c := client.NewClient(conf.APIKey, conf.SecretKey)
	c.Logger.SetOutput(io.Discard)

	// Return service
	return &Binance{
		Client: c,
	}, conf.Validate()
}
