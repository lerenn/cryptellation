package entities

type SymbolListener struct {
	Exchange    string `gorm:"primaryKey;autoIncrement:false"`
	Pair        string `gorm:"primaryKey;autoIncrement:false"`
	Subscribers int64
}
