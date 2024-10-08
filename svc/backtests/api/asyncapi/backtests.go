// Backtests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"
	"github.com/lerenn/cryptellation/pkg/utils"

	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"

	"github.com/google/uuid"
)

func (msg *CreateRequestMessage) Set(payload client.BacktestCreationPayload) {
	msg.Payload.StartTime = DateSchema(payload.StartTime)
	msg.Payload.EndTime = (*DateSchema)(payload.EndTime)
	msg.Payload.Accounts = accountModelsToAPI(payload.Accounts)
}

func (msg *GetRequestMessage) Set(backtestID uuid.UUID) {
	msg.Payload.Id = BacktestIDSchema(backtestID.String())
}

func (msg *GetResponseMessage) Set(backtest backtest.Backtest) {
	schema := BacktestSchema{
		Id:                  BacktestIDSchema(backtest.ID.String()),
		StartTime:           DateSchema(backtest.StartTime),
		EndTime:             DateSchema(backtest.EndTime),
		PeriodBetweenEvents: PeriodSchema(backtest.PeriodBetweenEvents),
	}

	msg.Payload.Backtest = &schema
}

func (msg *GetResponseMessage) ToModel() (backtest.Backtest, error) {
	p, err := period.FromString(string(msg.Payload.Backtest.PeriodBetweenEvents))
	if err != nil {
		return backtest.Backtest{}, err
	}

	return backtest.Backtest{
		ID:                  uuid.MustParse(string(msg.Payload.Backtest.Id)),
		StartTime:           time.Time(msg.Payload.Backtest.StartTime),
		EndTime:             time.Time(msg.Payload.Backtest.EndTime),
		PeriodBetweenEvents: p,
	}, nil
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

func (msg *ListRequestMessage) Set() error {
	// No payload
	return nil
}

func (msg *ListResponseMessage) Set(backtests []backtest.Backtest) {
	msg.Payload.Backtests = make([]BacktestSchema, len(backtests))
	for i, b := range backtests {
		msg.Payload.Backtests[i] = BacktestSchema{
			Id:                  BacktestIDSchema(b.ID.String()),
			StartTime:           DateSchema(b.StartTime),
			EndTime:             DateSchema(b.EndTime),
			PeriodBetweenEvents: PeriodSchema(b.PeriodBetweenEvents),
		}
	}
}

func (msg *ListResponseMessage) ToModel() ([]backtest.Backtest, error) {
	backtests := make([]backtest.Backtest, len(msg.Payload.Backtests))
	for i, b := range msg.Payload.Backtests {
		p, err := period.FromString(string(b.PeriodBetweenEvents))
		if err != nil {
			return nil, err
		}

		backtests[i] = backtest.Backtest{
			ID:                  uuid.MustParse(string(b.Id)),
			StartTime:           time.Time(b.StartTime),
			EndTime:             time.Time(b.EndTime),
			PeriodBetweenEvents: p,
		}
	}

	return backtests, nil
}
