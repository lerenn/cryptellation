package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	redisKeySymbolListenerPrefix = "symbol-listeners-count"
	redisKeySymbolListener       = redisKeySymbolListenerPrefix + "-%s-%s"
)

func (db *Adapter) IncrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	key := fmt.Sprintf(redisKeySymbolListener, exchange, pairSymbol)
	return db.redis.Client.Incr(ctx, key).Result()
}

func (db *Adapter) DecrementSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	key := fmt.Sprintf(redisKeySymbolListener, exchange, pairSymbol)
	return db.redis.Client.Decr(ctx, key).Result()
}

func (db *Adapter) GetSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) (int64, error) {
	key := fmt.Sprintf(redisKeySymbolListener, exchange, pairSymbol)
	content, err := db.redis.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(content, 10, 64)
}

func (db *Adapter) ClearAllSymbolListenersCount(ctx context.Context) error {
	keys, err := db.redis.Client.Keys(ctx, redisKeySymbolListenerPrefix+"*").Result()
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err := db.redis.Client.Set(ctx, k, 0, 0).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Adapter) ClearSymbolListenerSubscribers(ctx context.Context, exchange, pairSymbol string) error {
	key := fmt.Sprintf(redisKeySymbolListener, exchange, pairSymbol)
	_, err := db.redis.Client.Set(ctx, key, 0, time.Second).Result()
	return err
}
