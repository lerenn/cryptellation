package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
)

// Forwardtest is the entity for a forwardtest.
type Forwardtest struct {
	ID        string    `bson:"_id"`
	UpdatedAt time.Time `bson:"updated_at"`

	Accounts map[string]Account `bson:"accounts"`
	Orders   []Order            `bson:"orders"`
}

// ToModel converts a Forwardtest entity to a Forwardtest model.
func (ft Forwardtest) ToModel() (forwardtest.Forwardtest, error) {
	id, err := uuid.Parse(ft.ID)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	orders, err := ToOrderModels(ft.Orders)
	if err != nil {
		return forwardtest.Forwardtest{}, err
	}

	return forwardtest.Forwardtest{
		ID:        id,
		UpdatedAt: ft.UpdatedAt,
		Accounts:  ToAccountModels(ft.Accounts),
		Orders:    orders,
	}, nil
}

// FromForwardtestModel converts a Forwardtest model to a Forwardtest entity.
func FromForwardtestModel(ft forwardtest.Forwardtest) Forwardtest {
	return Forwardtest{
		ID:        ft.ID.String(),
		UpdatedAt: time.Now(),
		Accounts:  FromAccountModels(ft.Accounts),
		Orders:    FromOrderModels(ft.Orders),
	}
}
