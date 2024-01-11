package domain

import (
	"context"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/account"
)

func (b Backtests) GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}
