package migrations

import (
	"log"

	"gorm.io/gorm"
)

func InitSchema(tx *gorm.DB) error {
	log.Println("No migration detected: running the init schema.")

	type SymbolListener struct {
		Exchange    string `gorm:"primaryKey;autoIncrement:false"`
		PairSymbol  string `gorm:"primaryKey;autoIncrement:false"`
		Subscribers int64
	}

	return tx.AutoMigrate(&SymbolListener{})
}
