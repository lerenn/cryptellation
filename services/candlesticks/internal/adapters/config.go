package adapters

import (
	dbAdapters "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db"
	exchangesAdapter "github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/exchanges"
)

type Config struct {
	Exchanges exchangesAdapter.Config
	Database  dbAdapters.Config
}
