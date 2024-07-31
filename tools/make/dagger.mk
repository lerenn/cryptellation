dagger/develop:
	@dagger develop -m ./build/ci/dagger
	@dagger develop -m ./pkg/dagger
	@dagger develop -m ./svc/backtests/build/ci/dagger
	@dagger develop -m ./svc/backtests/pkg/dagger
	@dagger develop -m ./svc/candlesticks/build/ci/dagger
	@dagger develop -m ./svc/candlesticks/pkg/dagger
	@dagger develop -m ./svc/exchanges/build/ci/dagger
	@dagger develop -m ./svc/exchanges/pkg/dagger
	@dagger develop -m ./svc/forwardtests/build/ci/dagger
	@dagger develop -m ./svc/forwardtests/pkg/dagger
	@dagger develop -m ./svc/indicators/build/ci/dagger
	@dagger develop -m ./svc/indicators/pkg/dagger
	@dagger develop -m ./svc/ticks/build/ci/dagger
	@dagger develop -m ./svc/ticks/pkg/dagger