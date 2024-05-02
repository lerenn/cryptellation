package asyncapi

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/utils"
	client "github.com/lerenn/cryptellation/svc/backtests/clients/go"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
)

func (msg *OrdersListResponseMessage) Set(orders []order.Order) {
	msg.Payload.Orders = make([]OrderSchema, len(orders))
	for i, o := range orders {
		msg.Payload.Orders[i] = orderModelToAPI(o)
	}
}

func (msg *OrdersCreateRequestMessage) Set(payload client.OrderCreationPayload) {
	// Backtest
	msg.Payload.Id = BacktestIDSchema(payload.BacktestID)

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

	// Return order
	return order.Order{
		ID:            uint64(utils.FromReferenceOrDefault(o.Id)),
		ExecutionTime: (*time.Time)(o.ExecutionTime),
		Type:          t,
		Exchange:      string(o.Exchange),
		Pair:          string(o.Pair),
		Side:          s,
		Quantity:      o.Quantity,
		Price:         utils.FromReferenceOrDefault(o.Price),
	}, nil
}

func orderModelToAPI(o order.Order) OrderSchema {
	return OrderSchema{
		Id:            (*int64)(utils.ToReference((int64)(o.ID))),
		ExecutionTime: (*DateSchema)(o.ExecutionTime),
		Type:          OrderTypeSchema(o.Type.String()),
		Exchange:      ExchangeSchema(o.Exchange),
		Pair:          PairSchema(o.Pair),
		Side:          OrderSideSchema(o.Side.String()),
		Quantity:      o.Quantity,
		Price:         &o.Price,
	}
}
