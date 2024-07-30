package entities

import (
	"cryptellation/pkg/models/account"
)

type Balance struct {
	AssetName string  `bson:"asset_name"`
	Exchange  string  `bson:"exchange"`
	Balance   float64 `bson:"balance"`
}

func ToAccountModels(balances []Balance) map[string]account.Account {
	models := make(map[string]account.Account)
	for _, b := range balances {
		if _, exists := models[b.Exchange]; !exists {
			models[b.Exchange] = account.Account{
				Balances: make(map[string]float64),
			}
		}

		models[b.Exchange].Balances[b.AssetName] = b.Balance
	}
	return models
}

func FromAccountModels(accounts map[string]account.Account) []Balance {
	entities := make([]Balance, 0)

	for exchange, account := range accounts {
		for asset, balance := range account.Balances {
			entities = append(entities, Balance{
				AssetName: asset,
				Exchange:  exchange,
				Balance:   balance,
			})
		}
	}

	return entities
}
