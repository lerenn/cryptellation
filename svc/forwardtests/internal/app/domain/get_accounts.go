package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
)

func (ft *ForwardTests) GetAccounts(ctx context.Context, forwardTestID uuid.UUID) (map[string]account.Account, error) {
	f, err := ft.db.ReadForwardTest(ctx, forwardTestID)
	if err != nil {
		return nil, err
	}

	return f.Accounts, nil
}
