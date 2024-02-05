// Ticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g application -p asyncapi -i ../asyncapi.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g user        -p asyncapi -i ../asyncapi.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.30.2 -g types       -p asyncapi -i ../asyncapi.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	client "github.com/lerenn/cryptellation/svc/ticks/clients/go"
	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (msg *RegisteringRequestMessage) Set(payload client.TicksFilterPayload) {
	msg.Payload.Exchange = ExchangeSchema(payload.Exchange)
	msg.Payload.Pair = PairSchema(payload.Pair)
}

func (msg *TickMessage) ToModel() tick.Tick {
	return tick.Tick{
		Time:     time.Time(msg.Payload.Time),
		Pair:     string(msg.Payload.Pair),
		Price:    msg.Payload.Price,
		Exchange: string(msg.Payload.Exchange),
	}
}
