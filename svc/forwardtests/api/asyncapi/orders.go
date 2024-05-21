package asyncapi

import (
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/forwardtests/clients/go"
)

func (msg *OrdersCreateRequestMessage) Set(payload client.OrderCreationPayload) {
	// Backtest
	msg.Payload.Id = ForwardTestIDSchema(payload.ForwardTestID.String())

	// Order
	msg.Payload.Order.Exchange = ExchangeSchema(payload.Exchange)
	msg.Payload.Order.Pair = PairSchema(payload.Pair)
	msg.Payload.Order.Type = OrderTypeSchema(payload.Type)
	msg.Payload.Order.Side = OrderSideSchema(payload.Side)
	msg.Payload.Order.Quantity = payload.Quantity
}

func (msg OrdersCreateRequestMessage) ToModel() (order.Order, error) {
	return orderModelFromAPI(msg.Payload.Order)
}

func orderModelFromAPI(o OrderSchema) (order.Order, error) {
	// Check type
	t := order.Type(o.Type)
	if err := t.Validate(); err != nil {
		return order.Order{}, err
	}

	// Check side
	s := order.Side(o.Side)
	if err := s.Validate(); err != nil {
		return order.Order{}, err
	}

	// Check IF
	id, err := uuid.Parse(utils.FromReferenceOrDefault(o.Id))
	if err != nil && id != uuid.Nil {
		return order.Order{}, err
	}

	// Return order
	return order.Order{
		ID:            id,
		ExecutionTime: (*time.Time)(o.ExecutionTime),
		Type:          t,
		Exchange:      string(o.Exchange),
		Pair:          string(o.Pair),
		Side:          s,
		Quantity:      o.Quantity,
		Price:         utils.FromReferenceOrDefault(o.Price),
	}, nil
}
