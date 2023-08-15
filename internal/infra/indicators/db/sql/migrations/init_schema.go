package migrations

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	log.Println("No migration detected: running the init schema.")

	type SimpleMovingAverage struct {
		ExchangeName string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PairSymbol   string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PeriodSymbol string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PeriodNumber int       `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		PriceType    string    `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Time         time.Time `gorm:"primaryKey;autoIncrement:false;index:candlestick,unique"`
		Price        float64
	}

	return tx.AutoMigrate(&SimpleMovingAverage{})
}
