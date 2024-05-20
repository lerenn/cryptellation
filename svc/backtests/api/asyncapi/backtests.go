// Backtests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/event"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (msg *CreateRequestMessage) Set(payload client.BacktestCreationPayload) {
	msg.Payload.StartTime = DateSchema(payload.StartTime)
	msg.Payload.EndTime = (*DateSchema)(payload.EndTime)
	msg.Payload.Accounts = accountModelsToAPI(payload.Accounts)
}

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

func (msg *SubscribeRequestMessage) Set(backtestID uuid.UUID, exchange, pair string) {
	msg.Payload.Id = BacktestIDSchema(backtestID.String())
	msg.Payload.Exchange = ExchangeSchema(exchange)
	msg.Payload.Pair = PairSchema(pair)
}

func (msg *AdvanceRequestMessage) Set(backtestID uuid.UUID) {
	msg.Payload.Id = BacktestIDSchema(backtestID.String())
}

func (msg *CreateRequestMessage) ToModel() (backtest.NewPayload, error) {
	// Format accounts
	accounts := make(map[string]account.Account)
	for _, acc := range msg.Payload.Accounts {
		name, a := accountModelFromAPI(acc)
		accounts[name] = a
	}

	// Get duration between events
	var duration *time.Duration
	if msg.Payload.Period != nil {
		symbol, err := period.FromString((string)(*msg.Payload.Period))
		if err != nil {
			return backtest.NewPayload{}, err
		}

		duration = utils.ToReference(symbol.Duration())
	}

	// Return model
	return backtest.NewPayload{
		Accounts:              accounts,
		StartTime:             time.Time(msg.Payload.StartTime),
		EndTime:               (*time.Time)(msg.Payload.EndTime),
		DurationBetweenEvents: duration,
	}, nil
}

func (msg *EventMessage) Set(evt event.Event) error {
	msg.Payload.Time = DateSchema(evt.Time)
	msg.Payload.Type = evt.Type.String()

	// Set message depending on event type
	switch evt.Type {
	case event.TypeIsStatus:
		statusEvt, ok := evt.Content.(event.Status)
		if !ok {
			return event.ErrMismatchingType
		}

		msg.Payload.Content.Finished = statusEvt.Finished
	case event.TypeIsTick:
		t, ok := evt.Content.(tick.Tick)
		if !ok {
			return event.ErrMismatchingType
		}

		msg.Payload.Content.Exchange = ExchangeSchema(t.Exchange)
		msg.Payload.Content.Pair = PairSchema(t.Pair)
		msg.Payload.Content.Price = t.Price
		msg.Payload.Content.Time = DateSchema(t.Time)
	default:
		return event.ErrUnknownType
	}

	return nil
}
