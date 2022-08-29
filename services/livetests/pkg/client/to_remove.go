package client

import (
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/account"
	"github.com/digital-feather/cryptellation/services/livetests/pkg/client/proto"
)

func accountToProtoBuf(a account.Account) *proto.Account {
	assets := make(map[string]float64)
	for n, a := range a.Balances {
		assets[n] = float64(a)
	}

	return &proto.Account{
		Assets: assets,
	}
}
