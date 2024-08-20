package domain

import (
	"context"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"

	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app"
	"github.com/lerenn/cryptellation/svc/forwardtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

func (ft ForwardTests) List(ctx context.Context, _ app.ListFilters) ([]forwardtest.ForwardTest, error) {
	telemetry.L(ctx).Info("Listing forward tests")

	list, err := ft.db.ListForwardTests(ctx, db.ListFilters{})
	if err != nil {
		telemetry.L(ctx).Errorf("Error listing forward tests: %q", err.Error())
		return nil, err
	}

	telemetry.L(ctx).Infof("Listed %d forward tests", len(list))
	return list, nil
}
