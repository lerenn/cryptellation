package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrderQty = errors.New("invalid-order-quantity")
)

type Order struct {
	ID            uuid.UUID
	ExecutionTime *time.Time
	Type          Type
	Exchange      string
	Pair          string
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
