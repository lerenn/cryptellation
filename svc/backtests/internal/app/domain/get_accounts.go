package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/models/account"
)

func (b Backtests) GetAccounts(ctx context.Context, backtestId uuid.UUID) (map[string]account.Account, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}
