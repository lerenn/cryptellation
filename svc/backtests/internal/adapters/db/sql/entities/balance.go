package entities

import (
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
)

type Balance struct {
	AssetName    string `gorm:"primaryKey"`
	ExchangeName string `gorm:"primaryKey"`
	BacktestID   uint   `gorm:"primaryKey"`
	Balance      float64
}

func ToAccountModels(balances []Balance) map[string]account.Account {
	models := make(map[string]account.Account)
	for _, b := range balances {
		if _, exists := models[b.ExchangeName]; !exists {
			models[b.ExchangeName] = account.Account{
				Balances: make(map[string]float64),
			}
		}

		models[b.ExchangeName].Balances[b.AssetName] = b.Balance
	}
	return models
}

func FromAccountModels(backtestID uint, accounts map[string]account.Account) []Balance {
	entities := make([]Balance, 0)

	for exchange, account := range accounts {
		for asset, balance := range account.Balances {
			entities = append(entities, Balance{
				AssetName:    asset,
				BacktestID:   backtestID,
				ExchangeName: exchange,
				Balance:      balance,
			})
		}
	}

	return entities
}
