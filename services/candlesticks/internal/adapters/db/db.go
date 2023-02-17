package dbAdapters

import (
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/adapters/db/sql"
	"github.com/digital-feather/cryptellation/services/candlesticks/internal/components/candlesticks/db"
)

func New(c Config) (db.Port, error) {
	return sql.New(c.SQL)
}
