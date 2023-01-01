package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/digital-feather/cryptellation/services/backtests/internal/application/domains/backtest"
	"github.com/digital-feather/cryptellation/services/backtests/internal/application/ports/db"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

const (
	redisKeyBacktestIDs   = "backtests"
	redisKeyBacktest      = "backtest-%d"
	redisKeyMutexBacktest = "backtest-lock-%d"
)

var (
	lockOptions = []redsync.Option{
		redsync.WithExpiry(db.Expiration),
		redsync.WithRetryDelay(db.RetryDelay),
		redsync.WithTries(db.Tries),
	}
)

type DB struct {
	client     *redis.Client
	lockClient *redsync.Redsync
}

func New() (*DB, error) {
	var c Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading redis config: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	// TODO Check connection

	pool := goredis.NewPool(client)
	ls := redsync.New(pool)

	return &DB{
		client:     client,
		lockClient: ls,
	}, nil
}

func (d *DB) CreateBacktest(ctx context.Context, bt *backtest.Backtest) error {
	incr, err := d.client.Incr(ctx, redisKeyBacktestIDs).Result()
	if err != nil {
		return err
	}

	bt.ID = uint(incr)
	return d.client.Set(ctx, backtestKey(uint(incr)), bt, 0).Err()
}

func (d *DB) ReadBacktest(ctx context.Context, id uint) (backtest.Backtest, error) {
	bt := backtest.Backtest{}

	bValue, err := d.client.Get(ctx, backtestKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return bt, db.ErrRecordNotFound
	} else if err != nil {
		return bt, err
	}

	if err := json.Unmarshal(bValue, &bt); err != nil {
		return bt, err
	}

	return bt, nil
}

func (d *DB) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	return d.client.Set(ctx, backtestKey(bt.ID), bt, 0).Err()
}

func (d *DB) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	return d.client.Del(ctx, backtestKey(bt.ID)).Err()
}

func (d *DB) LockedBacktest(ctx context.Context, id uint, fn db.LockedBacktestCallback) error {
	mutex := d.lockClient.NewMutex(backtestMutexName(id), lockOptions...)
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

	bt, err := d.ReadBacktest(ctx, id)
	if err != nil {
		return err
	}

	err = fn(&bt)
	if err != nil {
		return err
	}

	return d.UpdateBacktest(ctx, bt)
}

func backtestKey(id uint) string {
	return fmt.Sprintf(redisKeyBacktest, id)
}

func backtestMutexName(id uint) string {
	return fmt.Sprintf(redisKeyMutexBacktest, id)
}
