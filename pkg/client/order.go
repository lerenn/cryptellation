package client

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/order"
)

type OrderCreationPayload struct {
	RunID    uuid.UUID
	Type     order.Type
	Exchange string
	Pair     string
	Side     order.Side
	Quantity float64
}
