.DEFAULT_GOAL  := help

DAGGER_COMMAND := _EXPERIMENTAL_DAGGER_INTERACTIVE_TUI=true dagger run go run ./build/ci/dagger.go

.PHONY: build
build: ## Run the build
	@${DAGGER_COMMAND} build

.PHONY: ci
ci: ## Run the CI
	@${DAGGER_COMMAND} all

.PHONY: lint
lint: ## Lint the code
	@${DAGGER_COMMAND} linter

.PHONY: generate
generate: ## Generate specified code across the codebase
	@${DAGGER_COMMAND} generator

.PHONY: test
test: ## Run tests
	@${DAGGER_COMMAND} test

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
