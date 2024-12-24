package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidOrderQty is the error for an invalid order quantity.
	ErrInvalidOrderQty = errors.New("invalid-order-quantity")
)

// Order is the struct for an order.
type Order struct {
	ID            uuid.UUID  `json:"id"`
	ExecutionTime *time.Time `json:"execution_time"`
	Type          Type       `json:"type"`
	Exchange      string     `json:"exchange"`
	Pair          string     `json:"pair"`
	Side          Side       `json:"side"`
	Quantity      float64    `json:"quantity"`
	Price         float64    `json:"price"`
}

// Validate validates the order.
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
