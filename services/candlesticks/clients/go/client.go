package client

import "github.com/digital-feather/cryptellation/services/candlesticks/clients/go/internal/controllers/nats"

func New(c Config) Interfacer {
	switch c.Type {
	case "nats":
		fallthrough
	default:
		return nats.New()
	}
}
