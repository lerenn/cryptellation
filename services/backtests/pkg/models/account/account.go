package account

import (
	"errors"

	"github.com/digital-feather/cryptellation/services/backtests/pkg/client/proto"
)

var (
	ErrInvalidBalanceAmount = errors.New("invalid-balance-amount")
	ErrInvalidBalanceAsset  = errors.New("invalid-balance-asset")
)

type Account struct {
	Balances map[string]float64
}

func (a Account) Validate() error {
	for asset, balance := range a.Balances {
		if asset == "" {
			return ErrInvalidBalanceAsset
		}

		if balance < 0 {
			return ErrInvalidBalanceAmount
		}
	}

	return nil
}

func (a Account) ToProtoBuf() *proto.Account {
	assets := make(map[string]float32)
	for n, a := range a.Balances {
		assets[n] = float32(a)
	}

	return &proto.Account{
		Assets: assets,
	}
}

func FromProtoBuf(pb *proto.Account) Account {
	assets := make(map[string]float64)
	for n, a := range pb.Assets {
		assets[n] = float64(a)
	}

	return Account{
		Balances: assets,
	}
}
