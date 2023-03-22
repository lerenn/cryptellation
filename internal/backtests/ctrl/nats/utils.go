package nats

import (
	"time"

	"github.com/digital-feather/cryptellation/internal/backtests/infra/events/nats/generated"
	"github.com/digital-feather/cryptellation/pkg/types/account"
	"github.com/digital-feather/cryptellation/pkg/types/order"
	"github.com/digital-feather/cryptellation/pkg/utils"
)

func accountModelToAPI(name string, account account.Account) generated.AccountSchema {
	assets := make([]generated.AssetSchema, 0, len(account.Balances))
	for name, qty := range account.Balances {
		assets = append(assets, generated.AssetSchema{
			Name:   name,
			Amount: qty,
		})
	}

	return generated.AccountSchema{
		Name:   name,
		Assets: assets,
	}
}

func accountModelFromAPI(a generated.AccountSchema) (string, account.Account) {
	assets := make(map[string]float64)
	for _, asset := range a.Assets {
		assets[asset.Name] = asset.Amount
	}

	return a.Name, account.Account{
		Balances: assets,
	}
}

func orderModelFromAPI(o generated.OrderSchema) (order.Order, error) {
	// Check type
	t := order.Type(o.Type)
	if err := t.Validate(); err != nil {
		return order.Order{}, err
	}

	// Check side
	s := order.Side(o.Side)
	if err := s.Validate(); err != nil {
		return order.Order{}, err
	}

	// Return order
	return order.Order{
		ID:            uint64(utils.FromReferenceOrDefault(o.ID)),
		ExecutionTime: (*time.Time)(o.ExecutionTime),
		Type:          t,
		ExchangeName:  string(o.ExchangeName),
		PairSymbol:    string(o.PairSymbol),
		Side:          s,
		Quantity:      o.Quantity,
		Price:         utils.FromReferenceOrDefault(o.Price),
	}, nil
}

func orderModelToAPI(o order.Order) generated.OrderSchema {
	return generated.OrderSchema{
		ID:            (*int64)(utils.ToReference((int64)(o.ID))),
		ExecutionTime: (*generated.DateSchema)(o.ExecutionTime),
		Type:          generated.OrderTypeSchema(o.Type.String()),
		ExchangeName:  generated.ExchangeNameSchema(o.ExchangeName),
		PairSymbol:    generated.PairSymbolSchema(o.PairSymbol),
		Side:          generated.OrderSideSchema(o.Side.String()),
		Quantity:      o.Quantity,
		Price:         &o.Price,
	}
}
