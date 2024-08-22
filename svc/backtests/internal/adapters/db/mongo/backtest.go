package backtests

import (
	"context"
	"fmt"

	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"

	"github.com/lerenn/cryptellation/svc/backtests/internal/adapters/db/mongo/entities"
	port "github.com/lerenn/cryptellation/svc/backtests/internal/app/ports/db"
	"github.com/lerenn/cryptellation/svc/backtests/pkg/backtest"

	"github.com/google/uuid"
)

const (
	// CollectionName is the name of the collection in the database
	CollectionName = "backtests"
)

func (a *Adapter) CreateBacktest(ctx context.Context, bt backtest.Backtest) error {
	// Check ID is not nil
	if bt.ID == uuid.Nil {
		return ErrNilID
	}

	// Create backtest
	entity := entities.FromBacktestModel(bt)
	_, err := a.client.Collection(CollectionName).InsertOne(ctx, entity)
	return err
}

func (a *Adapter) ReadBacktest(ctx context.Context, id uuid.UUID) (backtest.Backtest, error) {
	var entity entities.Backtest

	// Check ID is not nil
	if id == uuid.Nil {
		return backtest.Backtest{}, ErrNilID
	}

	// Get object from database
	err := a.client.
		Collection(CollectionName).
		FindOne(ctx, map[string]any{"_id": id.String()}).
		Decode(&entity)
	if err != nil {
		return backtest.Backtest{}, err
	}

	// Transform entity to model
	return entity.ToModel()
}

func (a *Adapter) ListBacktests(ctx context.Context) ([]backtest.Backtest, error) {
	var bts []backtest.Backtest

	// Get objects from database
	cursor, err := a.client.Collection(CollectionName).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Transform entities to models
	for cursor.Next(ctx) {
		var entity entities.Backtest
		if err := cursor.Decode(&entity); err != nil {
			return nil, err
		}

		bt, err := entity.ToModel()
		if err != nil {
			return nil, err
		}

		bts = append(bts, bt)
	}

	return bts, nil
}

func (a *Adapter) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	// Check ID is not nil
	if bt.ID == uuid.Nil {
		return ErrNilID
	}

	// Update backtest
	entity := entities.FromBacktestModel(bt)
	_, err := a.client.
		Collection(CollectionName).
		ReplaceOne(ctx, map[string]any{"_id": bt.ID.String()}, entity)
	return err
}

func (a *Adapter) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	// Check ID is not nil
	if bt.ID == uuid.Nil {
		return ErrNilID
	}

	// Delete backtest
	_, err := a.client.Collection(CollectionName).DeleteOne(ctx, map[string]any{"_id": bt.ID.String()})
	return err
}

func (a *Adapter) LockedBacktest(ctx context.Context, id uuid.UUID, fn port.LockedBacktestCallback) (err error) {
	var entity entities.Backtest

	// Check ID is not nil
	if id == uuid.Nil {
		return ErrNilID
	}

	// Get backtest from database
	err = a.client.
		Collection(CollectionName).
		FindOneAndUpdate(ctx,
			map[string]any{"_id": id.String()},
			map[string]any{
				"$set": map[string]any{
					"locked": map[string]any{
						"by":  "cryptellation-backtests",
						"key": uuid.New(),
					},
				},
			}).
		Decode(&entity)
	if err != nil {
		return err
	}

	// Transform entity to model
	bt, err := entity.ToModel()
	if err != nil {
		return err
	}

	// Recover from panic
	defer func() {
		if r := recover(); r != nil {
			telemetry.L(ctx).Info(fmt.Sprint("Recovered in f", r))
		}

		localErr := a.UpdateBacktest(ctx, bt)
		if localErr != nil {
			err = localErr
		}
	}()

	// Call callback
	if err := fn(&bt); err != nil {
		return err
	}

	// Update backtest in database
	return a.UpdateBacktest(ctx, bt)
}
