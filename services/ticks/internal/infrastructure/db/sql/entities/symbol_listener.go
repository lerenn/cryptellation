package entities

type SymbolListener struct {
	Exchange    string `gorm:"primaryKey;autoIncrement:false"`
	PairSymbol  string `gorm:"primaryKey;autoIncrement:false"`
	Subscribers int64
}
