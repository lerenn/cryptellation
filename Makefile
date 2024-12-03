DAGGER_CMD         := dagger call -m ./build/ci
DOCKER_COMPOSE_CMD := docker compose -f ./deployments/docker-compose/cryptellation.docker-compose.yaml
PROJECT_ROOT_PATH  := .

.DEFAULT_GOAL     := help

.PHONY: check
check: generate lint test ## Generate, lint and test the code

.PHONY: clean 
clean: local/down ## Clean the project
	@$(DOCKER_COMPOSE_CMD) rm
	@$(MAKE) -C deployments clean

.PHONY: dagger/check-generation
dagger/check-generation: ## Run all checks for generated code through Dagger
	@$(DAGGER_CMD) check-generation --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/develop
dagger/develop: ## Run Dagger develop on all Dagger modules
	@dagger develop -m ./pkg/dagger
	@dagger develop -m ./build/ci/dagger

.PHONY: dagger/lint
dagger/lint: ## Run all linters through Dagger
	@$(DAGGER_CMD) linter --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/tests/unit
dagger/tests/unit: ## Run all unit tests through Dagger
	@$(DAGGER_CMD) unit-tests --source-dir=$(PROJECT_ROOT_PATH) stdout
	
.PHONY: generate
generate: ## Generate the code
	@go generate ./...

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

.PHONY: lint
lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run ./...

.PHONY: local/down
local/down: ## Stop the local environment
	@$(DOCKER_COMPOSE_CMD) down

.PHONY: local/pull
local/pull: ## Pull the local environment images
	@$(DOCKER_COMPOSE_CMD) pull

.PHONY: local/up
local/up: ## Start the local environment
	@$(DOCKER_COMPOSE_CMD) up -d

.PHONY: test
test: test/unit test/integration test/end-to-end ## Launch all tests

.PHONY: test/unit
test/unit: ## Launch unit tests
	@go test $$(go list ./... | grep -v -e /activities -e /test)

.PHONY: test/integration
test/integration: local/up ## Launch integration tests
	@go test $$(go list ./pkg/domains/... | grep /activities)

.PHONY: test/end-to-end
test/end-to-end: local/up ## Launch end-to-end tests
	@go test ./test/...