package client

import (
	"context"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
)

type Client interface {
	CreateLivetest(ctx context.Context, accounts map[string]account.Account) (id uint, err error)
}
