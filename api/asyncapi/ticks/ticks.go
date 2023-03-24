// Backtests
//go:generate asyncapi-codegen -g application -p ticks -i ./../ticks.yaml -o ./app.gen.go
//go:generate asyncapi-codegen -g client      -p ticks -i ./../ticks.yaml -o ./client.gen.go
//go:generate asyncapi-codegen -g broker      -p ticks -i ./../ticks.yaml -o ./broker.gen.go
//go:generate asyncapi-codegen -g types       -p ticks -i ./../ticks.yaml -o ./types.gen.go
//go:generate asyncapi-codegen -g nats        -p ticks -i ./../ticks.yaml -o ./nats.gen.go

package ticks

import (
	"time"

	client "github.com/digital-feather/cryptellation/clients/go"
	"github.com/digital-feather/cryptellation/pkg/tick"
)

func (msg *RegisteringRequestMessage) Set(payload client.TicksFilterPayload) {
	msg.Payload.Exchange = ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = PairSymbolSchema(payload.PairSymbol)
}

func (msg *TickMessage) ToModel() tick.Tick {
	return tick.Tick{
		Time:       time.Time(msg.Payload.Time),
		PairSymbol: string(msg.Payload.PairSymbol),
		Price:      msg.Payload.Price,
		Exchange:   string(msg.Payload.Exchange),
	}
}
