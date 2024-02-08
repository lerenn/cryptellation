package migrations

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	telemetry.L(tx.Statement.Context).Info("No migration detected: running the init schema.")

	type Candlestick struct {
		Exchange   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Pair       string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Period     string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Time       time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Open       float64
		High       float64
		Low        float64
		Close      float64
		Volume     float64
		Uncomplete bool
	}

	return tx.AutoMigrate(&Candlestick{})
}
