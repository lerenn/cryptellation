package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
)

// Order is the entity for an order.
type Order struct {
	ID            string     `bson:"id"`
	ExecutionTime *time.Time `bson:"execution_time"`
	Type          string     `bson:"type"`
	Exchange      string     `bson:"exchange"`
	Pair          string     `bson:"pair"`
	Side          string     `bson:"side"`
	Quantity      float64    `bson:"quantity"`
	Price         float64    `bson:"price"`
}

// ToModel converts the entity to a model.
func (o Order) ToModel() (order.Order, error) {
	t := order.Type(o.Type)
	if err := t.Validate(); err != nil {
		return order.Order{}, err
	}

	s := order.Side(o.Side)
	if err := s.Validate(); err != nil {
		return order.Order{}, err
	}

	id, err := uuid.Parse(o.ID)
	if err != nil {
		return order.Order{}, err
	}

	return order.Order{
		ID:            id,
		ExecutionTime: o.ExecutionTime,
		Type:          t,
		Exchange:      o.Exchange,
		Pair:          o.Pair,
		Side:          s,
		Quantity:      o.Quantity,
		Price:         o.Price,
	}, nil
}

// ToOrderModels converts a slice of entities to a slice of models.
func ToOrderModels(orders []Order) ([]order.Order, error) {
	var err error
	models := make([]order.Order, len(orders))
	for i, e := range orders {
		if models[i], err = e.ToModel(); err != nil {
			return nil, err
		}
	}
	return models, nil
}

// FromOrderModels converts a slice of models into a slice of entities.
func FromOrderModels(models []order.Order) []Order {
	entities := make([]Order, len(models))
	for i, m := range models {
		entities[i] = FromOrderModel(m)
	}
	return entities
}

// FromOrderModel converts a model into an entity.
func FromOrderModel(m order.Order) Order {
	return Order{
		ID:            m.ID.String(),
		ExecutionTime: m.ExecutionTime,
		Type:          m.Type.String(),
		Exchange:      m.Exchange,
		Pair:          m.Pair,
		Side:          m.Side.String(),
		Quantity:      m.Quantity,
		Price:         m.Price,
	}
}
