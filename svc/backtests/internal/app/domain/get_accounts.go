package domain

import (
	"context"

	"cryptellation/pkg/models/account"

	"github.com/google/uuid"
)

func (b Backtests) GetAccounts(ctx context.Context, backtestId uuid.UUID) (map[string]account.Account, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}
