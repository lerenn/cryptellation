package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/v1/api"
	"github.com/lerenn/cryptellation/v1/pkg/models/account"
	"github.com/lerenn/cryptellation/v1/pkg/models/forwardtest"
	"github.com/lerenn/cryptellation/v1/pkg/models/order"
)

type Forwardtest struct {
	ID uuid.UUID

	cryptellation client
}

func (ft Forwardtest) CreateOrder(
	ctx context.Context,
	order order.Order,
) (api.CreateForwardtestOrderWorkflowResults, error) {
	return ft.cryptellation.Client.CreateForwardtestOrder(ctx, api.CreateForwardtestOrderWorkflowParams{
		ForwardtestID: ft.ID,
		Order:         order,
	})
}

func (ft Forwardtest) ListAccounts(
	ctx context.Context,
) (map[string]account.Account, error) {
	res, err := ft.cryptellation.Client.ListForwardtestAccounts(ctx, api.ListForwardtestAccountsWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return nil, err
	}

	return res.Accounts, nil
}

func (ft Forwardtest) GetStatus(
	ctx context.Context,
) (forwardtest.Status, error) {
	res, err := ft.cryptellation.Client.GetForwardtestStatus(ctx, api.GetForwardtestStatusWorkflowParams{
		ForwardtestID: ft.ID,
	})
	if err != nil {
		return forwardtest.Status{}, err
	}

	return res.Status, nil
}
