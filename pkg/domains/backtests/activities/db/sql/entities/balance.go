package entities

import (
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
)

// Balance is the entity for a balance.
type Balance struct {
	AssetName string  `json:"asset_name"`
	Exchange  string  `json:"exchange"`
	Balance   float64 `json:"balance"`
}

// ToAccountModels transforms account entities to account models.
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

// FromAccountModels transforms account models to account entities.
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
