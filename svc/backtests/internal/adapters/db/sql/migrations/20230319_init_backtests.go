package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

const migration20230319Title = "20230219-init-backtests"

var migration20230319 = gormigrate.Migration{
	ID: migration20230319Title,
	Migrate: func(tx *gorm.DB) error {
		telemetry.L(tx.Statement.Context).Info("Running migration" + migration20230319Title)

		type Balance struct {
			AssetName  string `gorm:"primaryKey"`
			Exchange   string `gorm:"primaryKey"`
			BacktestID uint   `gorm:"primaryKey"`
			Balance    float64
		}

		type Order struct {
			ID            uint `gorm:"primaryKey"`
			BacktestID    uint
			ExecutionTime *time.Time
			Type          string
			Exchange      string
			Pair          string
			Side          string
			Quantity      float64
			Price         float64
		}

		type TickSubscription struct {
			ID         uint `gorm:"primaryKey"`
			BacktestID uint
			Exchange   string
			Pair       string
		}

		type Backtest struct {
			ID                  uint `gorm:"primaryKey"`
			StartTime           time.Time
			CurrentTime         time.Time
			CurrentPriceType    string
			EndTime             time.Time
			PeriodBetweenEvents string
			Balances            []Balance          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
			Orders              []Order            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
			TickSubscriptions   []TickSubscription `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		}

		return tx.AutoMigrate(&Backtest{}, &Order{}, &Balance{}, &TickSubscription{})
	},
	Rollback: func(tx *gorm.DB) error {
		telemetry.L(tx.Statement.Context).Info("Running rollback: " + migration20230319Title)
		return tx.Migrator().DropTable("tick_subscriptions", "balances", "orders", "backtests")
	},
}
