package backtests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	port "github.com/lerenn/cryptellation/internal/components/backtests/ports/db"
	"github.com/lerenn/cryptellation/pkg/models/backtest"
)

const (
	redisKeyBacktestIDs   = "backtests"
	redisKeyBacktest      = "backtest-%d"
	redisKeyMutexBacktest = "backtest-lock-%d"
)

var (
	lockOptions = []redsync.Option{
		redsync.WithExpiry(port.Expiration),
		redsync.WithRetryDelay(port.RetryDelay),
		redsync.WithTries(port.Tries),
	}
)

func (db *Adapter) CreateBacktest(ctx context.Context, bt *backtest.Backtest) error {
	incr, err := db.redis.Client.Incr(ctx, redisKeyBacktestIDs).Result()
	if err != nil {
		return err
	}

	bt.ID = uint(incr)
	return db.redis.Client.Set(ctx, backtestKey(uint(incr)), bt, 0).Err()
}

func (db *Adapter) ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error) {
	bt := backtest.Backtest{}

	bValue, err := db.redis.Client.Get(ctx, backtestKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return bt, port.ErrRecordNotFound
	} else if err != nil {
		return bt, err
	}

	if err := json.Unmarshal(bValue, &bt); err != nil {
		return bt, err
	}

	return bt, nil
}

func (db *Adapter) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	return db.redis.Client.Set(ctx, backtestKey(bt.ID), bt, 0).Err()
}

func (db *Adapter) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	return db.redis.Client.Del(ctx, backtestKey(bt.ID)).Err()
}

func (db *Adapter) LockedBacktest(ctx context.Context, id uint, fn port.LockedBacktestCallback) error {
	mutex := db.redis.ClientLock.NewMutex(backtestMutexName(id), lockOptions...)
	if err := mutex.Lock(); err != nil {
		return err
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}

		ok, localErr := mutex.Unlock()
		if localErr != nil {
			err = localErr
		} else if !ok {
			err = fmt.Errorf("unlock failed for backtest %d", id)
		}
	}()

	bt, err := db.ReadBacktest(ctx, id)
	if err != nil {
		return err
	}

	err = fn(&bt)
	if err != nil {
		return err
	}

	return db.UpdateBacktest(ctx, bt)
}

func backtestKey(id uint) string {
	return fmt.Sprintf(redisKeyBacktest, id)
}

func backtestMutexName(id uint) string {
	return fmt.Sprintf(redisKeyMutexBacktest, id)
}
