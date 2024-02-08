package migrations

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	telemetry.L(tx.Statement.Context).Info("No migration detected: running the init schema.")

	type SimpleMovingAverage struct {
		Exchange     string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Pair         string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Period       string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PeriodNumber int       `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PriceType    string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Time         time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Price        float64
	}

	return tx.AutoMigrate(&SimpleMovingAverage{})
}
