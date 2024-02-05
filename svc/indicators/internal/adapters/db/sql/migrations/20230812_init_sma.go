package migrations

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const migration20230812Title = "20230812-init-sma"

var migration20230812 = gormigrate.Migration{
	ID: migration20230812Title,
	Migrate: func(tx *gorm.DB) error {
		log.Println("Running migration:", migration20230812Title)

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
	},
	Rollback: func(tx *gorm.DB) error {
		log.Println("Running rollback:", migration20230812Title)
		return tx.Migrator().DropTable("simple_moving_averages")
	},
}
