package entities

import (
	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

type ForwardTest struct {
	ID       string             `bson:"_id"`
	Accounts map[string]Account `bson:"accounts"`
}

func (f ForwardTest) ToModel() (forwardtest.ForwardTest, error) {
	id, err := uuid.Parse(f.ID)
	if err != nil {
		return forwardtest.ForwardTest{}, err
	}

	return forwardtest.ForwardTest{
		ID:       id,
		Accounts: ToAccountModels(f.Accounts),
	}, nil
}

func FromForwardTestModel(f forwardtest.ForwardTest) ForwardTest {
	return ForwardTest{
		ID:       f.ID.String(),
		Accounts: FromAccountModels(f.Accounts),
	}
}
