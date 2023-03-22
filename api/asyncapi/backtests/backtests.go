// Backtests
//go:generate asyncapi-codegen -g application -p backtests -i ./../backtests.yaml -o ./app.gen.go
//go:generate asyncapi-codegen -g client      -p backtests -i ./../backtests.yaml -o ./client.gen.go
//go:generate asyncapi-codegen -g broker      -p backtests -i ./../backtests.yaml -o ./broker.gen.go
//go:generate asyncapi-codegen -g types       -p backtests -i ./../backtests.yaml -o ./types.gen.go
//go:generate asyncapi-codegen -g nats        -p backtests -i ./../backtests.yaml -o ./nats.gen.go

package backtests

import (
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/internal/backtests/app/domain"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/event"
	"github.com/digital-feather/cryptellation/pkg/types/period"
	"github.com/digital-feather/cryptellation/pkg/types/tick"
	"github.com/digital-feather/cryptellation/pkg/utils"
)

func (msg *BacktestsCreateRequestMessage) Set(payload client.BacktestCreationPayload) {
	msg.Payload.StartTime = DateSchema(payload.StartTime)
	msg.Payload.EndTime = (*DateSchema)(payload.EndTime)
	// TODO
}

// func accountModelsToAPI(accounts map[string]account.Account) []generated.AccountSchema {
// 	apiAccounts := make([]generated.AccountSchema, 0, len(accounts))
// 	for name, acc := range accounts {
// 		apiAccounts = append(apiAccounts, generated.AccountSchema{
// 			Name: name,
// 		})
// 	}
// }

func (msg *BacktestsCreateRequestMessage) ToModel() (domain.NewPayload, error) {
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
			return domain.NewPayload{}, err
		}

		duration = utils.ToReference(symbol.Duration())
	}

	// Return model
	return domain.NewPayload{
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
