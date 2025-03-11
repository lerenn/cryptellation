package entities

import "github.com/lerenn/cryptellation/v1/pkg/models/account"

// Account is the entity for an account.
type Account struct {
	Balances map[string]float64 `json:"balances"`
}

// ToAccountModels converts a map of Account to a map of account.Account.
func ToAccountModels(accounts map[string]Account) map[string]account.Account {
	models := make(map[string]account.Account)
	for exchange, acc := range accounts {
		models[exchange] = account.Account{
			Balances: acc.Balances,
		}
	}
	return models
}

// FromAccountModels converts a map of account.Account to a map of Account.
func FromAccountModels(accounts map[string]account.Account) map[string]Account {
	entities := make(map[string]Account)

	for exchange, acc := range accounts {
		entities[exchange] = Account{
			Balances: acc.Balances,
		}
	}

	return entities
}
