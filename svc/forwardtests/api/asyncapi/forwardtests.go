// ForwardTests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../pkg/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../pkg/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../pkg/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import (
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

func (msg *CreateRequestMessage) Set(payload forwardtest.NewPayload) {
	// Format accounts
	accounts := make([]AccountSchema, 0, len(payload.Accounts))
	for name, acc := range payload.Accounts {
		accounts = append(accounts, accountModelToAPI(name, acc))
	}

	// Set accounts
	msg.Payload.Accounts = accounts
}

func (msg CreateRequestMessage) ToModel() (forwardtest.NewPayload, error) {
	// Format accounts
	accounts := make(map[string]account.Account)
	for _, acc := range msg.Payload.Accounts {
		name, a := accountModelFromAPI(acc)
		accounts[name] = a
	}

	// Return model
	return forwardtest.NewPayload{
		Accounts: accounts,
	}, nil
}
