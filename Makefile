.PHONY: all
.DEFAULT_GOAL := help

SHORT_COMMIT_SHA := $(shell git rev-parse --short HEAD)

DOCKER_IMAGE_TAG=$(SHORT_COMMIT_SHA)
DOCKER_COMPOSE := docker-compose -p cryptellation

SERVICES := $(patsubst services/%, %, $(wildcard services/*))

.PHONY: docker/build
docker/build: $(addprefix docker/build/,$(SERVICES)) ## Build docker image

.PHONY: docker/build/%
docker/build/%:
	@docker build \
		-t digital-feather/cryptellation-$*:$(SHORT_COMMIT_SHA) \
		-f ./build/package/Dockerfile \
		--build-arg BINARY=$* \
		./

.PHONY: docker/push
docker/push: $(addprefix docker/push/,$(SERVICES)) ## Push docker image

.PHONY: docker/push/%
docker/push/%:
	@git diff-index --quiet HEAD || (echo "ERROR: Some files have been modified. Please commit before pushing."; exit 1)
	@docker push -t digital-feather/cryptellation-$*:$(SHORT_COMMIT_SHA)

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

.PHONY: test/clean
test/clean:
	$(MAKE) -C test/end-to-end clean
	@$(DOCKER_COMPOSE) -f ./test/integration/docker-compose.yml down
	@rm -f cover.out

.PHONY: test/end-to-end
test/end-to-end: ## Perform end-to-end tests
	@$(MAKE) -C test/end-to-end run

.PHONY: test/integration
test/integration: ## Perform integration tests
	@$(DOCKER_COMPOSE) -f ./test/integration/docker-compose.yml build
	@$(DOCKER_COMPOSE) -f ./test/integration/docker-compose.yml run tests

.PHONY: test/unit
test/unit: ## Perform unit tests
	@go test $(shell go list ./cmd/... ./pkg/... ./services/... | grep -v -e /io/) -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: test
test: test/unit test/integration test/end-to-end ## Run tests

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
