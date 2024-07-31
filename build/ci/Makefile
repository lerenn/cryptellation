include ../../tools/make/help.mk

.PHONY: lint
lint:
	@dagger call lint --source-dir=../..:source-only stdout

.PHONY: check-generation
check-generation:
	@dagger call check-generation --source-dir=../.. stdout