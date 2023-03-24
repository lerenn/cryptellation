package migrations

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const migration20230219Title = "20230219-init-candlesticks"

var migration20230219 = gormigrate.Migration{
	ID: migration20230219Title,
	Migrate: func(tx *gorm.DB) error {
		log.Println("Running migration:", migration20230219Title)

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
	},
	Rollback: func(tx *gorm.DB) error {
		log.Println("Running rollback:", migration20230219Title)
		return tx.Migrator().DropTable("candlesticks")
	},
}
