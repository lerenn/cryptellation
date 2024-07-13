package cache

import "time"

type option func(c *cache)

func WithMaxSize(maxSize int) option {
	return func(c *cache) {
		c.settings.maxSize = maxSize
	}
}

func WithExpirationTime(expirationTime time.Duration) option {
	return func(c *cache) {
		c.settings.expirationTime = expirationTime
	}
}
