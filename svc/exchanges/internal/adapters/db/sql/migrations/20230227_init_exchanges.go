package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

const migration20230227Title = "20230219-init-exchanges"

var migration20230227 = gormigrate.Migration{
	ID: migration20230227Title,
	Migrate: func(tx *gorm.DB) error {
		telemetry.L(tx.Statement.Context).Info("Running migration: " + migration20230227Title)

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
	},
	Rollback: func(tx *gorm.DB) error {
		telemetry.L(tx.Statement.Context).Info("Running rollback: " + migration20230227Title)
		return tx.Migrator().DropTable("exchanges", "pairs", "periods")
	},
}
