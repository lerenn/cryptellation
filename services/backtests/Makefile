.PHONY: all
.DEFAULT_GOAL := help

SERVICE_NAME = backtests

SHORT_COMMIT_SHA := $(shell git rev-parse --short HEAD)

DOCKER_IMAGE_TAG=$(SHORT_COMMIT_SHA)
DOCKER_COMPOSE := docker-compose -p cryptellation-$(SERVICE_NAME)

.PHONY: docker/build
docker/build: ## Build docker image
	@docker build -t digital-feather/cryptellation-$(SERVICE_NAME):$(SHORT_COMMIT_SHA) -f ./build/package/Dockerfile .

.PHONY: docker/push
docker/push: ## Push docker image
	@git diff-index --quiet HEAD || (echo "ERROR: Some files have been modified. Please commit before pushing."; exit 1)
	@docker push -t digital-feather/cryptellation-$(SERVICE_NAME):$(SHORT_COMMIT_SHA)

.PHONY: clean
clean: test/clean ## Clean everything

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

.PHONY: test/clean
test/clean:
	@$(DOCKER_COMPOSE) -f ./test/docker-compose.yml down
	@rm -f cover.out

.PHONY: test/end-to-end
test/end-to-end: ## Perform end-to-end tests
	@echo "TODO"

.PHONY: test/integration
test/integration: ## Perform integration tests
	@$(DOCKER_COMPOSE) -f ./test/docker-compose.yml build $(SERVICE_NAME)
	@$(DOCKER_COMPOSE) -f ./test/docker-compose.yml run $(SERVICE_NAME) ./test/integration.sh

.PHONY: test/unit
test/unit: ## Perform unit tests
	@go test $(shell go list ./... | grep -v -e /internal/infrastructure) -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: test
test: test/unit test/integration test/end-to-end ## Run tests

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
