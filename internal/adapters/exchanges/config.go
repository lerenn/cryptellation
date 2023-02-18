package exchanges

import (
	"github.com/digital-feather/cryptellation/internal/adapters/exchanges/binance"
)

type Config struct {
	Binance binance.Config
}

func LoadConfigFromEnv() (c Config) {
	c.Binance = binance.LoadConfigFromEnv()
	return c
}
