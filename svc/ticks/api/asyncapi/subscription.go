package asyncapi

import (
	"github.com/lerenn/cryptellation/pkg/models/event"
)

func (msg *ListeningNotificationMessage) FromModel(sub event.TickSubscription) {
	msg.Payload.Exchange = ExchangeSchema(sub.Exchange)
	msg.Payload.Pair = PairSchema(sub.Pair)
}

func (msg ListeningNotificationMessage) ToModel() event.TickSubscription {
	return event.TickSubscription{
		Exchange: string(msg.Payload.Exchange),
		Pair:     string(msg.Payload.Pair),
	}
}
