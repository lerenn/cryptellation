package migrations

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const migration20230303Title = "20230219-init-ticks"

var migration20230303 = gormigrate.Migration{
	ID: migration20230303Title,
	Migrate: func(tx *gorm.DB) error {
		log.Println("Running migration:", migration20230303Title)

		type SymbolListener struct {
			Exchange    string `gorm:"primaryKey;autoIncrement:false"`
			PairSymbol  string `gorm:"primaryKey;autoIncrement:false"`
			Subscribers int64
		}

		return tx.AutoMigrate(&SymbolListener{})
	},
	Rollback: func(tx *gorm.DB) error {
		log.Println("Running rollback:", migration20230303Title)
		return tx.Migrator().DropTable("symbol_listeners")
	},
}
