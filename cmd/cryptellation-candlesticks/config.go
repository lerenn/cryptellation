package main

import (
	"os"
	"strconv"

	"github.com/digital-feather/cryptellation/services/candlesticks"
)

func loadConfigFromEnv() (config candlesticks.ServiceConfig) {
	config.Adapters.Exchanges.Binance.ApiKey = os.Getenv("BINANCE_API_KEY")
	config.Adapters.Exchanges.Binance.SecretKey = os.Getenv("BINANCE_SECRET_KEY")

	config.Adapters.Database.SQL.Host = os.Getenv("SQLDB_HOST")
	config.Adapters.Database.SQL.Port, _ = strconv.Atoi(os.Getenv("SQLDB_PORT"))
	config.Adapters.Database.SQL.User = os.Getenv("SQLDB_USER")
	config.Adapters.Database.SQL.Password = os.Getenv("SQLDB_PASSWORD")
	config.Adapters.Database.SQL.Database = os.Getenv("SQLDB_DATABASE")

	config.Controllers.NATS.Host = os.Getenv("NATS_HOST")
	config.Controllers.NATS.Port, _ = strconv.Atoi(os.Getenv("NATS_PORT"))

	return config
}
