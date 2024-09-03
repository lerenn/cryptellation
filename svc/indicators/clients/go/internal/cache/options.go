package cache

type option func(c *cache)

func WithMaxSize(maxSize int) option {
	return func(c *cache) {
		c.settings.maxSize = maxSize
	}
}

func WithPreLoadingAfterSize(preLoadingAfterSize int) option {
	return func(c *cache) {
		c.settings.preLoadingAfterSize = preLoadingAfterSize
	}
}

func WithPreLoadingBeforeSize(preLoadingBeforeSize int) option {
	return func(c *cache) {
		c.settings.preLoadingBeforeSize = preLoadingBeforeSize
	}
}
