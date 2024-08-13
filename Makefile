PROJECT_ROOT_PATH := .

DAGGER_MODULE_PATH := ./build/ci

include $(PROJECT_ROOT_PATH)/tools/make/dagger.mk
include $(PROJECT_ROOT_PATH)/tools/make/help.mk
include $(PROJECT_ROOT_PATH)/tools/make/kind.mk
include $(PROJECT_ROOT_PATH)/tools/make/root.mk

.PHONY: dagger/develop
dagger/develop: ## Run Dagger develop on all Dagger modules
	@dagger develop -m ./internal/dagger
	@dagger develop -m ./svc/candlesticks/pkg/dagger
	@dagger develop -m ./svc/candlesticks/build/ci/dagger
	@dagger develop -m ./svc/ticks/pkg/dagger
	@dagger develop -m ./svc/ticks/build/ci/dagger
	@dagger develop -m ./svc/exchanges/pkg/dagger
	@dagger develop -m ./svc/exchanges/build/ci/dagger
	@dagger develop -m ./svc/backtests/pkg/dagger
	@dagger develop -m ./svc/backtests/build/ci/dagger
	@dagger develop -m ./svc/forwardtests/pkg/dagger
	@dagger develop -m ./svc/forwardtests/build/ci/dagger
	@dagger develop -m ./svc/indicators/pkg/dagger
	@dagger develop -m ./svc/indicators/build/ci/dagger
	@dagger develop -m ./build/ci/dagger