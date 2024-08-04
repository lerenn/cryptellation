PROJECT_ROOT_PATH := .

DAGGER_MODULE_PATH := ./build/ci

include tools/make/dagger.mk
include tools/make/help.mk
include tools/make/kind.mk
include tools/make/root.mk

.PHONY: dagger/develop
dagger/develop: ## Run Dagger develop on all Dagger modules
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
	@dagger develop -m ./build/ci/dagger