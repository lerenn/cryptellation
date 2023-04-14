.PHONY: all
.DEFAULT_GOAL := help

.PHONY: clean
clean: test/clean ## Clean everything
	$(MAKE) -C tools/minikube-env clean

.PHONY: lint
lint: ## Lint the code
ifeq ($(shell which golangci-lint &> /dev/null; echo $$?),1)
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif
	@LOG_LEVEL=error golangci-lint run

.PHONY: generate ## Generate specified code across the codebase
generate:
	@go generate ./...

.PHONY: test
test: ## Run tests
	@make -C test all

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
