package migrations

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	telemetry.L(tx.Statement.Context).Info("No migration detected: running the init schema.")

	type Period struct {
		Symbol string `gorm:"primaryKey;autoIncrement:false"`
	}

	type Pair struct {
		Symbol string `gorm:"primaryKey;autoIncrement:false"`
	}

	type Exchange struct {
		Name         string   `gorm:"primaryKey;autoIncrement:false"`
		Pairs        []Pair   `gorm:"many2many:exchanges_pairs;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Periods      []Period `gorm:"many2many:exchanges_periods;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Fees         float64
		LastSyncTime time.Time
	}

	return tx.AutoMigrate(&Period{}, &Pair{}, &Exchange{})
}
