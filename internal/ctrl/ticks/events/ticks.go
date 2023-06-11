// Ticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g application -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g client      -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./client.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g broker      -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./broker.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g types       -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./types.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.13.1 -g nats        -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./nats.gen.go

package events

import (
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/tick"
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
