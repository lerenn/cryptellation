package asyncapi

import "github.com/lerenn/cryptellation/pkg/models/account"

func accountModelFromAPI(a AccountSchema) (string, account.Account) {
	assets := make(map[string]float64)
	for _, asset := range a.Assets {
		assets[asset.Name] = asset.Amount
	}

	return a.Name, account.Account{
		Balances: assets,
	}
}

func accountModelToAPI(name string, account account.Account) AccountSchema {
	assets := make([]AssetSchema, 0, len(account.Balances))
	for name, qty := range account.Balances {
		assets = append(assets, AssetSchema{
			Name:   name,
			Amount: qty,
		})
	}

	return AccountSchema{
		Name:   name,
		Assets: assets,
	}
}
