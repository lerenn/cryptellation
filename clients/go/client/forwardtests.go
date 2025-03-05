package client

import (
	"context"

	"github.com/lerenn/cryptellation/v1/api"
)

func (c client) NewForwardtest(
	ctx context.Context,
	params api.CreateForwardtestWorkflowParams,
) (Forwardtest, error) {
	res, err := c.Client.CreateForwardtest(ctx, params)
	return Forwardtest{
		ID:            res.ID,
		cryptellation: c,
	}, err
}

func (c client) ListForwardtests(
	ctx context.Context,
	params api.ListForwardtestsWorkflowParams,
) ([]Forwardtest, error) {
	res, err := c.Client.ListForwardtests(ctx, params)
	if err != nil {
		return nil, err
	}

	forwardtests := make([]Forwardtest, len(res.Forwardtests))
	for i, ft := range res.Forwardtests {
		forwardtests[i] = Forwardtest{
			ID:            ft.ID,
			cryptellation: c,
		}
	}

	return forwardtests, nil
}
