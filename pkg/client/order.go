package client

import (
	"cryptellation/pkg/models/order"

	"github.com/google/uuid"
)

type OrderCreationPayload struct {
	RunID    uuid.UUID
	Type     order.Type
	Exchange string
	Pair     string
	Side     order.Side
	Quantity float64
}
