package controllers

import "github.com/digital-feather/cryptellation/services/candlesticks/internal/controllers/nats"

type Config struct {
	NATS nats.Config
}
