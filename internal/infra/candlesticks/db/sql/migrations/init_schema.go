package migrations

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	log.Println("No migration detected: running the init schema.")

	type Candlestick struct {
		ExchangeName string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PairSymbol   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PeriodSymbol string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Time         time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Open         float64
		High         float64
		Low          float64
		Close        float64
		Volume       float64
		Uncomplete   bool
	}

	return tx.AutoMigrate(&Candlestick{})
}
