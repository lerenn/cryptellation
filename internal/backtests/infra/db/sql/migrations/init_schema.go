package migrations

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	log.Println("No migration detected: running the init schema.")

	type Balance struct {
		AssetName    string `gorm:"primaryKey"`
		ExchangeName string `gorm:"primaryKey"`
		BacktestID   uint   `gorm:"primaryKey"`
		Balance      float64
	}

	type Order struct {
		ID            uint `gorm:"primaryKey"`
		BacktestID    uint
		ExecutionTime *time.Time
		Type          string
		ExchangeName  string
		PairSymbol    string
		Side          string
		Quantity      float64
		Price         float64
	}

	type TickSubscription struct {
		ID           uint `gorm:"primaryKey"`
		BacktestID   uint
		ExchangeName string
		PairSymbol   string
	}

	type Backtest struct {
		ID                  uint `gorm:"primaryKey"`
		StartTime           time.Time
		CurrentTime         time.Time
		CurrentPriceType    string
		EndTime             time.Time
		PeriodBetweenEvents string
		Balances            []Balance
		Orders              []Order
		TickSubscriptions   []TickSubscription
	}

	return tx.AutoMigrate(&Backtest{}, &Order{}, &Balance{}, &TickSubscription{})
}
