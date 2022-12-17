.PHONY: all
.DEFAULT_GOAL := help

DOCKER_COMPOSE := docker-compose -p cryptellation-exchanges

.PHONY: clean
clean: test/clean ## Clean everything
	@$(DOCKER_COMPOSE) down
	@rm -f cover.out

.PHONY: proto
proto: ## Generate protobuf code
	@mkdir -p pkg/client/proto
	@for PROTOFILE in $(shell find . -name *.proto); do \
		echo "Generating $$PROTOFILE"; \
		protoc \
			--proto_path=./api \
			--go_out=./pkg/client/proto \
			--go_opt=paths=source_relative \
			--go-grpc_opt=require_unimplemented_servers=false \
			--go-grpc_out=./pkg/client/proto \
			--go-grpc_opt=paths=source_relative \
			"$$PROTOFILE" || exit $?; \
	done

.PHONY: lint
lint: ## Lint the code
ifeq ($(shell which golangci-lint &> /dev/null; echo $$?),1)
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif
	@LOG_LEVEL=error golangci-lint run

.PHONY: generate ## Generate specified code across the codebase
generate:
	@go generate ./...

.PHONY: test/end-to-end
test/end-to-end: ## Perform end-to-end tests
	@echo "TODO"

.PHONY: test/integration
test/integration: ## Perform integration tests
	@$(DOCKER_COMPOSE) -f ./tests/integration/docker-compose.yml build exchanges
	@$(DOCKER_COMPOSE) -f ./tests/integration/docker-compose.yml run exchanges

.PHONY: test/unit
test/unit: ## Perform unit tests
	@go test $(shell go list ./... | grep -v -e /internal/adapters) -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: test
test: test/unit test/integration test/end-to-end ## Run tests

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
