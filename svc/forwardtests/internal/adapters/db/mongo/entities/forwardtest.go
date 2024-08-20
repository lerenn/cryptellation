package entities

import (
	"time"

	"cryptellation/svc/forwardtests/pkg/forwardtest"

	"github.com/google/uuid"
)

type ForwardTest struct {
	ID        string    `bson:"_id"`
	UpdatedAt time.Time `bson:"updated_at"`

	Accounts map[string]Account `bson:"accounts"`
	Orders   []Order            `bson:"orders"`
}

func (ft ForwardTest) ToModel() (forwardtest.ForwardTest, error) {
	id, err := uuid.Parse(ft.ID)
	if err != nil {
		return forwardtest.ForwardTest{}, err
	}

	orders, err := ToOrderModels(ft.Orders)
	if err != nil {
		return forwardtest.ForwardTest{}, err
	}

	return forwardtest.ForwardTest{
		ID:        id,
		UpdatedAt: ft.UpdatedAt,
		Accounts:  ToAccountModels(ft.Accounts),
		Orders:    orders,
	}, nil
}

func FromForwardTestModel(ft forwardtest.ForwardTest) ForwardTest {
	return ForwardTest{
		ID:        ft.ID.String(),
		UpdatedAt: time.Now(),
		Accounts:  FromAccountModels(ft.Accounts),
		Orders:    FromOrderModels(ft.Orders),
	}
}
