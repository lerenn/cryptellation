//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g application -p asyncapi -i ../asyncapi.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g user        -p asyncapi -i ../asyncapi.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.39.0 -g types       -p asyncapi -i ../asyncapi.yaml -o ./types.gen.go

package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
)

func (msg *TickMessage) FromModel(t tick.Tick) {
	msg.Payload.Exchange = ExchangeSchema(t.Exchange)
	msg.Payload.Pair = PairSchema(t.Pair)
	msg.Payload.Price = t.Price
	msg.Payload.Time = DateSchema(t.Time)
}

func (msg TickMessage) ToModel() tick.Tick {
	return tick.Tick{
		Time:     time.Time(msg.Payload.Time),
		Pair:     string(msg.Payload.Pair),
		Price:    msg.Payload.Price,
		Exchange: string(msg.Payload.Exchange),
	}
}
