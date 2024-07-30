package domain

import (
	"context"

	"cryptellation/pkg/models/account"

	"github.com/google/uuid"
)

func (ft *ForwardTests) GetAccounts(ctx context.Context, forwardTestID uuid.UUID) (map[string]account.Account, error) {
	f, err := ft.db.ReadForwardTest(ctx, forwardTestID)
	if err != nil {
		return nil, err
	}

	return f.Accounts, nil
}
