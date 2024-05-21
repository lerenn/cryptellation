package entities

import "github.com/lerenn/cryptellation/pkg/models/account"

type Account struct {
	Balances map[string]float64 `bson:"balances"`
}

func ToAccountModels(accounts map[string]Account) map[string]account.Account {
	models := make(map[string]account.Account)
	for exchange, acc := range accounts {
		models[exchange] = account.Account{
			Balances: acc.Balances,
		}
	}
	return models
}

func FromAccountModels(accounts map[string]account.Account) map[string]Account {
	entities := make(map[string]Account)

	for exchange, acc := range accounts {
		entities[exchange] = Account{
			Balances: acc.Balances,
		}
	}

	return entities
}
