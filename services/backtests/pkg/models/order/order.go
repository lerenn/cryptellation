package order

import (
	"errors"
	"time"

	"github.com/digital-feather/cryptellation/services/backtests/clients/go/proto"
)

var (
	ErrInvalidOrderQty = errors.New("invalid-order-quantity")
)

type Order struct {
	ID            uint64
	ExecutionTime *time.Time
	Type          Type
	ExchangeName  string
	PairSymbol    string
	Side          Side
	Quantity      float64
	Price         float64
}

func (o Order) Validate() error {
	if err := o.Type.Validate(); err != nil {
		return err
	}

	if err := o.Side.Validate(); err != nil {
		return err
	}

	if o.Quantity <= 0 {
		return ErrInvalidOrderQty
	}

	return nil
}

func (o Order) ToProtoBuf() *proto.Order {
	var pt *string
	if o.ExecutionTime != nil {
		t := o.ExecutionTime.UTC().Format(time.RFC3339Nano)
		pt = &t
	}

	return &proto.Order{
		Id:            o.ID,
		ExecutionTime: pt,
		Type:          o.Type.String(),
		ExchangeName:  o.ExchangeName,
		PairSymbol:    o.PairSymbol,
		Side:          o.Side.String(),
		Quantity:      float64(o.Quantity),
		Price:         float64(o.Price),
	}
}

func FromProtoBuf(pb *proto.Order) (Order, error) {
	var pt *time.Time
	if pb.ExecutionTime != nil {
		t, err := time.Parse(time.RFC3339Nano, *pb.ExecutionTime)
		if err != nil {
			return Order{}, err
		}
		pt = &t
	}

	ty := Type(pb.Type)
	if err := ty.Validate(); err != nil {
		return Order{}, err
	}

	s := Side(pb.Side)
	if err := s.Validate(); err != nil {
		return Order{}, err
	}

	return Order{
		ID:            pb.Id,
		ExecutionTime: pt,
		Type:          ty,
		ExchangeName:  pb.ExchangeName,
		PairSymbol:    pb.PairSymbol,
		Side:          s,
		Quantity:      float64(pb.Quantity),
		Price:         float64(pb.Price),
	}, nil
}
