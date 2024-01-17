package ci

import "dagger.io/dagger"

func DefaultExchangesVariables(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			WithEnvVariable("BINANCE_API_KEY", "").
			WithEnvVariable("BINANCE_SECRET_KEY", "")
	}
}

// BinanceDependency returns a function that set variables to use binance as a service
func BinanceDependency(client *dagger.Client) func(r *dagger.Container) *dagger.Container {
	return func(r *dagger.Container) *dagger.Container {
		return r.
			With(Secret(client, "BINANCE_API_KEY")).
			With(Secret(client, "BINANCE_SECRET_KEY"))
	}
}
