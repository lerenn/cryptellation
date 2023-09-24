// Backtests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g application -p events -i ./../../../../api/asyncapi/backtests.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g user        -p events -i ./../../../../api/asyncapi/backtests.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g types       -p events -i ./../../../../api/asyncapi/backtests.yaml -o ./types.gen.go

package events

import (
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/backtest"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/models/period"
	"github.com/lerenn/cryptellation/pkg/models/tick"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func (msg *BacktestsCreateRequestMessage) Set(payload client.BacktestCreationPayload) {
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

func (msg *BacktestsSubscribeRequestMessage) Set(backtestID uint, exchange, pair string) {
	msg.Payload.ID = BacktestIDSchema(backtestID)
	msg.Payload.ExchangeName = ExchangeNameSchema(exchange)
	msg.Payload.PairSymbol = PairSymbolSchema(pair)
}

func (msg *BacktestsAdvanceRequestMessage) Set(backtestID uint) {
	msg.Payload.ID = BacktestIDSchema(backtestID)
}

func (msg *BacktestsCreateRequestMessage) ToModel() (backtest.NewPayload, error) {
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

func (msg *BacktestsEventMessage) Set(evt event.Event) error {
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

		msg.Payload.Content.Exchange = ExchangeNameSchema(t.Exchange)
		msg.Payload.Content.PairSymbol = PairSymbolSchema(t.PairSymbol)
		msg.Payload.Content.Price = t.Price
		msg.Payload.Content.Time = DateSchema(t.Time)
	default:
		return event.ErrUnknownType
	}

	return nil
}
