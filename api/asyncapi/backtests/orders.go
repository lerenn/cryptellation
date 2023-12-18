package backtests

import (
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/order"
	"github.com/lerenn/cryptellation/pkg/utils"
)

func (msg *ListBacktestOrdersResponseMessage) Set(orders []order.Order) {
	msg.Payload.Orders = make([]OrderSchema, len(orders))
	for i, o := range orders {
		msg.Payload.Orders[i] = orderModelToAPI(o)
	}
}

func (msg *CreateBacktestOrderRequestMessage) Set(payload client.OrderCreationPayload) {
	// Backtest
	msg.Payload.Id = BacktestIDSchema(payload.BacktestID)

	// Order
	msg.Payload.Order.ExchangeName = ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Order.PairSymbol = PairSymbolSchema(payload.PairSymbol)
	msg.Payload.Order.Type = OrderTypeSchema(payload.Type)
	msg.Payload.Order.Side = OrderSideSchema(payload.Side)
	msg.Payload.Order.Quantity = payload.Quantity
}

func (msg CreateBacktestOrderRequestMessage) ToModel() (order.Order, error) {
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
		ExchangeName:  string(o.ExchangeName),
		PairSymbol:    string(o.PairSymbol),
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
		ExchangeName:  ExchangeNameSchema(o.ExchangeName),
		PairSymbol:    PairSymbolSchema(o.PairSymbol),
		Side:          OrderSideSchema(o.Side.String()),
		Quantity:      o.Quantity,
		Price:         &o.Price,
	}
}
