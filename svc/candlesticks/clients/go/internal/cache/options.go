package cache

type option func(c *cache)

func WithMaxSize(maxSize int) option {
	return func(c *cache) {
		c.settings.maxSize = maxSize
	}
}

func WithPreLoadingSize(preLoadingSize int) option {
	return func(c *cache) {
		c.settings.preLoadingSize = preLoadingSize
	}
}

func WithPreemptiveAsyncLoadingEnabled(preemptiveAsyncLoadingEnabled bool) option {
	return func(c *cache) {
		c.settings.preemptiveAsyncLoadingEnabled = preemptiveAsyncLoadingEnabled
	}
}
