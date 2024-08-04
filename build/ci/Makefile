include ../../tools/make/help.mk

.PHONY: check
check:
	@dagger call check --source-dir=../.. stdout

.PHONY: lint
lint:
	@dagger call lint --source-dir=../..:source-only stdout

.PHONY: check-generation
check-generation:
	@dagger call check-generation --source-dir=../.. stdout

.PHONY: unit-tests
unit-tests:
	@dagger call unit-tests --source-dir=../.. stdout

.PHONY: integration-tests
integration-tests:
	@dagger call integration-tests \
		--source-dir=../.. \
		--secrets-file=file:../../.credentials.env \
		stdout