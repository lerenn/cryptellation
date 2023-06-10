.PHONY: all
.DEFAULT_GOAL := help

.PHONY: all
all: unit integration end-to-end ## Run tests

.PHONY: clean
clean: ## Clean up tests
	@$(MAKE) -C integration clean
	@$(MAKE) -C end-to-end clean

.PHONY: unit
unit: ## Perform unit tests
	@go test $(shell go list ../cmd/... ../pkg/... ../internal/core/...) -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: integration
integration: ## Perform integration tests
	@$(MAKE) -C integration run

.PHONY: end-to-end
end-to-end: ## Perform end-to-end tests
	@$(MAKE) -C end-to-end run/local

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
