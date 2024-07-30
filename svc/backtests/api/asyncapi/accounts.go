package asyncapi

import (
	"cryptellation/pkg/models/account"

	"github.com/google/uuid"
)

func accountModelsToAPI(accounts map[string]account.Account) []AccountSchema {
	apiAccounts := make([]AccountSchema, 0, len(accounts))
	for accName, acc := range accounts {
		// Set assets
		assets := make([]AssetSchema, 0, len(acc.Balances))
		for assetName, amount := range acc.Balances {
			assets = append(assets, AssetSchema{
				Name:   assetName,
				Amount: amount,
			})
		}

		// Set account
		apiAccounts = append(apiAccounts, AccountSchema{
			Name:   accName,
			Assets: assets,
		})
	}

	return apiAccounts
}

func (msg *AccountsListRequestMessage) Set(backtestID uuid.UUID) {
	msg.Payload.Id = BacktestIDSchema(backtestID.String())
}

func (msg *AccountsListResponseMessage) Set(accounts map[string]account.Account) {
	// Format accounts
	respAccounts := make([]AccountSchema, 0, len(accounts))
	for name, acc := range accounts {
		respAccounts = append(respAccounts, accountModelToAPI(name, acc))
	}

	// Set response
	msg.Payload.Accounts = respAccounts
}

func (msg *AccountsListResponseMessage) ToModel() map[string]account.Account {
	accounts := make(map[string]account.Account)
	for _, accAPI := range msg.Payload.Accounts {
		name, accModel := accountModelFromAPI(accAPI)
		accounts[name] = accModel
	}
	return accounts
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

func accountModelFromAPI(a AccountSchema) (string, account.Account) {
	assets := make(map[string]float64)
	for _, asset := range a.Assets {
		assets[asset.Name] = asset.Amount
	}

	return a.Name, account.Account{
		Balances: assets,
	}
}
