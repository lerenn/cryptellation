// Exchanges
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml,../../../../internal/asyncapi/models.yaml -o ./types.gen.go

package asyncapi

import "github.com/lerenn/cryptellation/exchanges/pkg/exchange"

func (msg *ListRequestMessage) Set(names ...string) {
	msg.Payload = make([]ExchangeNameSchema, 0, len(names))
	for _, name := range names {
		msg.Payload = append(msg.Payload, ExchangeNameSchema(name))
	}
}

func (msg *ListRequestMessage) ToModel() []string {
	exchangesNames := make([]string, len(msg.Payload))
	for i, e := range msg.Payload {
		exchangesNames[i] = string(e)
	}
	return exchangesNames
}

func (msg *ListResponseMessage) Set(exchanges []exchange.Exchange) {
	msg.Payload.Exchanges = make([]ExchangeSchema, len(exchanges))
	for i, exch := range exchanges {
		// Periods
		periods := make([]PeriodSchema, len(exch.Periods))
		for j, p := range exch.Periods {
			periods[j] = PeriodSchema(p)
		}

		// Pairs
		pairs := make([]PairSchema, len(exch.Pairs))
		for j, p := range exch.Pairs {
			pairs[j] = PairSchema(p)
		}

		// Exchange
		msg.Payload.Exchanges[i] = ExchangeSchema{
			Fees:         exch.Fees,
			Name:         ExchangeNameSchema(exch.Name),
			Pairs:        pairs,
			Periods:      periods,
			LastSyncTime: exch.LastSyncTime,
		}
	}
}

func (msg *ListResponseMessage) ToModel() []exchange.Exchange {
	exchanges := make([]exchange.Exchange, len(msg.Payload.Exchanges))
	for i, exch := range msg.Payload.Exchanges {
		// Periods
		periods := make([]string, len(exch.Periods))
		for j, p := range exch.Periods {
			periods[j] = string(p)
		}

		// Pairs
		pairs := make([]string, len(exch.Pairs))
		for j, p := range exch.Pairs {
			pairs[j] = string(p)
		}

		exchanges[i] = exchange.Exchange{
			Name:         string(exch.Name),
			Fees:         exch.Fees,
			Pairs:        pairs,
			Periods:      periods,
			LastSyncTime: exch.LastSyncTime,
		}
	}

	return exchanges
}
