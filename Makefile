.PHONY: all
.DEFAULT_GOAL := help

.PHONY: clean
clean: ## Clean everything
	$(MAKE) -C deployments clean
	$(MAKE) -C test clean
	$(MAKE) -C build/package clean

.PHONY: lint
lint: ## Lint the code
ifeq ($(shell which golangci-lint &> /dev/null; echo $$?),1)
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif
	@LOG_LEVEL=error golangci-lint run

.PHONY: generate
generate: ## Generate specified code across the codebase
	@go generate ./...

.PHONY: test
test: ## Run tests
	@make -C test unit integration

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
