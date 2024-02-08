package migrations

import (
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	telemetry.L(tx.Statement.Context).Info("No migration detected: running the init schema.")

	type SymbolListener struct {
		Exchange    string `gorm:"primaryKey;autoIncrement:false"`
		Pair        string `gorm:"primaryKey;autoIncrement:false"`
		Subscribers int64
	}

	return tx.AutoMigrate(&SymbolListener{})
}
