.DEFAULT_GOAL := help

.PHONY: ci
ci: ## Execute all basic CI steps
	@dagger run go run ./tools/cmd/cryptellation-ci

.PHONY: generate
generate: ## Generate code files
	@dagger run go run ./tools/cmd/cryptellation-ci generate

.PHONY: lint
lint: ## Lint code
	@dagger run go run ./tools/cmd/cryptellation-ci lint

.PHONY: serve 
serve: ## Serve the Cryptellation stack for development
	@dagger run go run ./tools/cmd/cryptellation-ci serve

.PHONY: test
test: test/unit test/integration test/end-to-end ## Launch tests 

.PHONY: test/unit
test/unit: ## Launch unit tests
	@dagger run go run ./tools/cmd/cryptellation-ci test --type=unit

.PHONY: test/integration
test/integration: ## Launch integration tests
	@dagger run go run ./tools/cmd/cryptellation-ci test --type=integration

.PHONY: test/end-to-end
test/end-to-end: ## Launch end-to-end tests
	@dagger run go run ./tools/cmd/cryptellation-ci test --type=end-to-end

.PHONY: update
update: ## Update the dependencies
	@dagger run go run ./tools/cmd/cryptellation-ci update

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'