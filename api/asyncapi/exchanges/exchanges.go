// Backtests
//go:generate asyncapi-codegen -g application -p exchanges -i ./../exchanges.yaml -o ./app.gen.go
//go:generate asyncapi-codegen -g client      -p exchanges -i ./../exchanges.yaml -o ./client.gen.go
//go:generate asyncapi-codegen -g broker      -p exchanges -i ./../exchanges.yaml -o ./broker.gen.go
//go:generate asyncapi-codegen -g types       -p exchanges -i ./../exchanges.yaml -o ./types.gen.go
//go:generate asyncapi-codegen -g nats        -p exchanges -i ./../exchanges.yaml -o ./nats.gen.go

package exchanges

import "github.com/digital-feather/cryptellation/pkg/exchange"

func (msg *ExchangesRequestMessage) Set(names ...string) {
	msg.Payload = make([]ExchangeNameSchema, 0, len(names))
	for _, name := range names {
		msg.Payload = append(msg.Payload, ExchangeNameSchema(name))
	}
}

func (msg *ExchangesRequestMessage) ToModel() []string {
	exchangesNames := make([]string, len(msg.Payload))
	for i, e := range msg.Payload {
		exchangesNames[i] = string(e)
	}
	return exchangesNames
}

func (msg *ExchangesResponseMessage) Set(exchanges []exchange.Exchange) {
	msg.Payload.Exchanges = make([]ExchangeSchema, len(exchanges))
	for i, exch := range exchanges {
		// Periods
		periods := make([]PeriodSymbolSchema, len(exch.PeriodsSymbols))
		for j, p := range exch.PeriodsSymbols {
			periods[j] = PeriodSymbolSchema(p)
		}

		// Pairs
		pairs := make([]PairSymbolSchema, len(exch.PairsSymbols))
		for j, p := range exch.PairsSymbols {
			pairs[j] = PairSymbolSchema(p)
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

func (msg *ExchangesResponseMessage) ToModel() []exchange.Exchange {
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
			Name:           string(exch.Name),
			Fees:           exch.Fees,
			PairsSymbols:   pairs,
			PeriodsSymbols: periods,
			LastSyncTime:   exch.LastSyncTime,
		}
	}

	return exchanges
}
