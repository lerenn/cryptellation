// ForwardTests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"

	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
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

func (msg *ListResponseMessage) Set(payload []forwardtest.ForwardTest) {
	// Format forward tests
	tests := make([]ForwardTestSchema, 0, len(payload))
	for _, test := range payload {
		tests = append(tests, ForwardTestSchema{
			Id:        ForwardTestIDSchema(test.ID.String()),
			UpdatedAt: DateSchema(test.UpdatedAt),
		})
	}

	// Set forward tests
	msg.Payload.Forwardtests = tests
}

func (msg ListResponseMessage) ToModel() ([]forwardtest.ForwardTest, error) {
	// Format forward tests
	tests := make([]forwardtest.ForwardTest, 0, len(msg.Payload.Forwardtests))
	for _, ft := range msg.Payload.Forwardtests {
		id, err := uuid.Parse(string(ft.Id))
		if err != nil {
			return nil, err
		}

		tests = append(tests, forwardtest.ForwardTest{
			ID:        id,
			UpdatedAt: time.Time(ft.UpdatedAt),
		})
	}

	// Return forward tests
	return tests, nil
}
