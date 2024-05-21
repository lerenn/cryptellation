package forwardtest

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
)

var (
	ErrEmptyAccounts = errors.New("empty accounts")
)

type ForwardTest struct {
	ID       uuid.UUID
	Accounts map[string]account.Account
}

type NewPayload struct {
	Accounts map[string]account.Account
}

func (np NewPayload) Validate() error {
	if len(np.Accounts) == 0 {
		return ErrEmptyAccounts
	}

	return nil
}

func New(payload NewPayload) ForwardTest {
	return ForwardTest{
		ID:       uuid.New(),
		Accounts: payload.Accounts,
	}
}
