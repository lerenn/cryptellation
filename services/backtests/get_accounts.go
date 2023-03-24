package backtests

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/account"
)

func (b Backtests) GetAccounts(ctx context.Context, backtestId uint) (map[string]account.Account, error) {
	bt, err := b.db.ReadBacktest(ctx, backtestId)
	if err != nil {
		return nil, err
	}

	return bt.Accounts, nil
}
