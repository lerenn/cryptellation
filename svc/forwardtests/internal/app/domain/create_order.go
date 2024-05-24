package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/models/order"
	candlesticks "github.com/lerenn/cryptellation/svc/candlesticks/clients/go"
	"github.com/lerenn/cryptellation/svc/candlesticks/pkg/period"
)

func (f ForwardTests) CreateOrder(ctx context.Context, forwardTestID uuid.UUID, order order.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}

	telemetry.L(ctx).Debugf("Creating order %+v on forward test %q", order, forwardTestID.String())

	ft, err := f.db.ReadForwardTest(ctx, forwardTestID)
	if err != nil {
		return fmt.Errorf("could not get forward test from db: %w", err)
	}

	now := time.Now()
	list, err := f.candlesticks.Read(ctx, candlesticks.ReadCandlesticksPayload{
		Exchange: order.Exchange,
		Pair:     order.Pair,
		Period:   period.M1,
		Start:    &now,
		End:      &now,
		Limit:    1,
	})
	if err != nil {
		return fmt.Errorf("could not get candlesticks from service: %w", err)
	}

	_, cs, notEmpty := list.First()
	if !notEmpty {
		return fmt.Errorf("no data for order validation")
	}

	telemetry.L(ctx).Infof("Adding %+v order to %q forwardtest", order, forwardTestID.String())
	if err := ft.AddOrder(order, cs); err != nil {
		return err
	}

	return f.db.UpdateForwardTest(ctx, ft)
}
