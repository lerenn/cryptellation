.DEFAULT_GOAL := help

DAGGER_CMD    := dagger run go run ./tools/cmd/cryptellation-ci

.PHONY: ci
ci: ## Execute all basic CI steps
	@$(DAGGER_CMD)

.PHONY: generate
generate: ## Generate code files
	@$(DAGGER_CMD) generate

.PHONY: lint
lint: ## Lint code
	@$(DAGGER_CMD) lint

.PHONY: serve 
serve: ## Serve the Cryptellation stack for development
	@$(DAGGER_CMD) serve

.PHONY: test
test: test/unit test/integration test/end-to-end ## Launch tests 

.PHONY: test/unit
test/unit: ## Launch unit tests
	@$(DAGGER_CMD) test --type=unit

.PHONY: test/integration
test/integration: ## Launch integration tests
	@$(DAGGER_CMD) test --type=integration

.PHONY: test/end-to-end
test/end-to-end: ## Launch end-to-end tests
	@$(DAGGER_CMD) test --type=end-to-end

.PHONY: update
update: ## Update the dependencies
	@$(DAGGER_CMD) update

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'