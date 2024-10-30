// Backtests
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/pkg/models/event"

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
	msg.Payload.Mode = (*ModeSchema)(payload.Mode)
	msg.Payload.PricePeriod = (*PeriodSchema)(payload.PricePeriod)
}

func (msg *GetRequestMessage) Set(backtestID uuid.UUID) {
	msg.Payload.Id = BacktestIDSchema(backtestID.String())
}

func (msg *GetResponseMessage) Set(backtest backtest.Backtest) {
	tsub := make([]PricesSubscriptionSchema, len(backtest.PricesSubscriptions))
	for i, ts := range backtest.PricesSubscriptions {
		tsub[i] = PricesSubscriptionSchema{
			Exchange: ExchangeSchema(ts.Exchange),
			Pair:     PairSchema(ts.Pair),
		}
	}

	schema := BacktestSchema{
		Id: BacktestIDSchema(backtest.ID.String()),
		Parameters: BacktestParametersSchema{
			StartTime:   DateSchema(backtest.Parameters.StartTime),
			EndTime:     DateSchema(backtest.Parameters.EndTime),
			Mode:        ModeSchema(backtest.Parameters.Mode),
			PricePeriod: PeriodSchema(backtest.Parameters.PricePeriod),
		},
		PricesSubscriptions: tsub,
	}

	msg.Payload.Backtest = &schema
}

func (msg *GetResponseMessage) ToModel() (backtest.Backtest, error) {
	p, err := period.FromString(string(msg.Payload.Backtest.Parameters.PricePeriod))
	if err != nil {
		return backtest.Backtest{}, err
	}

	ts := make([]event.PricesSubscription, len(msg.Payload.Backtest.PricesSubscriptions))
	for i, t := range msg.Payload.Backtest.PricesSubscriptions {
		ts[i] = event.PricesSubscription{
			Exchange: string(t.Exchange),
			Pair:     string(t.Pair),
		}
	}

	return backtest.Backtest{
		ID: uuid.MustParse(string(msg.Payload.Backtest.Id)),
		Parameters: backtest.Parameters{
			StartTime:   time.Time(msg.Payload.Backtest.Parameters.StartTime),
			EndTime:     time.Time(msg.Payload.Backtest.Parameters.EndTime),
			Mode:        backtest.Mode(msg.Payload.Backtest.Parameters.Mode),
			PricePeriod: p,
		},
		PricesSubscriptions: ts,
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

	// Get period between prices
	var per *period.Symbol
	if msg.Payload.PricePeriod != nil {
		s, err := period.FromString((string)(*msg.Payload.PricePeriod))
		if err != nil {
			return backtest.NewPayload{}, err
		}
		per = &s
	}

	// Return model
	return backtest.NewPayload{
		Accounts:    accounts,
		StartTime:   time.Time(msg.Payload.StartTime),
		EndTime:     (*time.Time)(msg.Payload.EndTime),
		Mode:        (*backtest.Mode)(msg.Payload.Mode),
		PricePeriod: per,
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
	case event.TypeIsPrice:
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
			Id: BacktestIDSchema(b.ID.String()),
			Parameters: BacktestParametersSchema{
				StartTime:   DateSchema(b.Parameters.StartTime),
				EndTime:     DateSchema(b.Parameters.EndTime),
				PricePeriod: PeriodSchema(b.Parameters.PricePeriod),
			},
		}
	}
}

func (msg *ListResponseMessage) ToModel() ([]backtest.Backtest, error) {
	backtests := make([]backtest.Backtest, len(msg.Payload.Backtests))
	for i, b := range msg.Payload.Backtests {
		p, err := period.FromString(string(b.Parameters.PricePeriod))
		if err != nil {
			return nil, err
		}

		backtests[i] = backtest.Backtest{
			ID: uuid.MustParse(string(b.Id)),
			Parameters: backtest.Parameters{
				StartTime:   time.Time(b.Parameters.StartTime),
				EndTime:     time.Time(b.Parameters.EndTime),
				PricePeriod: p,
			},
		}
	}

	return backtests, nil
}
