.DEFAULT_GOAL := help

.PHONY: ci
ci: ## Execute all basic CI steps
	@dagger run go run ./cmd/ci

.PHONY: generate
generate: ## Generate code files
	@dagger run go run ./cmd/ci generate

.PHONY: lint
lint: ## Lint code
	@dagger run go run ./cmd/ci lint

.PHONY: serve 
serve: ## Serve the Cryptellation stack for development
	@dagger run go run ./cmd/ci serve

.PHONY: update
update: ## Update the dependencies
	@dagger run go run ./cmd/ci update

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'