DOCKER_COMPOSE_CMD := docker compose -f ./deployments/docker-compose/cryptellation.docker-compose.yaml 

.PHONE: all
all: generate lint test ## Generate, lint and test the code

.PHONY: clean
clean: local/down ## Clean the project
	@$(MAKE) -C ./clients/go clean
	@$(MAKE) -C ./cmd/cryptellation clean
	@$(MAKE) -C ./cmd/cryptellation-tui clean
	@$(MAKE) -C ./examples/go clean
	@$(MAKE) -C ./pkg clean
	@$(MAKE) -C ./svc/backtests clean
	@$(MAKE) -C ./svc/candlesticks clean
	@$(MAKE) -C ./svc/exchanges clean
	@$(MAKE) -C ./svc/forwardtests clean
	@$(MAKE) -C ./svc/indicators clean
	@$(MAKE) -C ./svc/ticks clean

.PHONY: build
build: docker/build ## Build the project

.PHONY: docker/build
docker/build: ## Build the docker images
	@$(MAKE) -C ./svc/backtests docker/build
	@$(MAKE) -C ./svc/candlesticks docker/build
	@$(MAKE) -C ./svc/exchanges docker/build
	@$(MAKE) -C ./svc/forwardtests docker/build
	@$(MAKE) -C ./svc/indicators docker/build
	@$(MAKE) -C ./svc/ticks docker/build

.PHONY: docker/publish
docker/publish: ## Publish the docker images
	@$(MAKE) -C ./svc/backtests docker/publish
	@$(MAKE) -C ./svc/candlesticks docker/publish
	@$(MAKE) -C ./svc/exchanges docker/publish
	@$(MAKE) -C ./svc/forwardtests docker/publish
	@$(MAKE) -C ./svc/indicators docker/publish
	@$(MAKE) -C ./svc/ticks docker/publish

.PHONY: local/down
local/down: ## Stop the local environment
	@$(DOCKER_COMPOSE_CMD) down

.PHONY: local/pull
local/pull: ## Pull the local environment images
	@$(DOCKER_COMPOSE_CMD) pull

.PHONY: local/up
local/up: ## Start the local environment
	@$(DOCKER_COMPOSE_CMD) up -d

.PHONY: go/generate
go/generate: ## Generate the golang code
	@$(MAKE) -C ./clients/go go/generate
	@$(MAKE) -C ./cmd/cryptellation go/generate
	@$(MAKE) -C ./cmd/cryptellation-tui go/generate
	@$(MAKE) -C ./examples/go go/generate
	@$(MAKE) -C ./pkg go/generate
	@$(MAKE) -C ./svc/backtests go/generate
	@$(MAKE) -C ./svc/candlesticks go/generate
	@$(MAKE) -C ./svc/exchanges go/generate
	@$(MAKE) -C ./svc/forwardtests go/generate
	@$(MAKE) -C ./svc/indicators go/generate
	@$(MAKE) -C ./svc/ticks go/generate

.PHONY: generate
generate: go/generate ## Generate the code

.PHONY: go/lint
go/lint: ## Lint the code
	@$(MAKE) -C ./clients/go go/lint
	@$(MAKE) -C ./cmd/cryptellation go/lint
	@$(MAKE) -C ./cmd/cryptellation-tui go/lint
	@$(MAKE) -C ./examples/go go/lint
	@$(MAKE) -C ./pkg go/lint
	@$(MAKE) -C ./svc/backtests go/lint
	@$(MAKE) -C ./svc/candlesticks go/lint
	@$(MAKE) -C ./svc/exchanges go/lint
	@$(MAKE) -C ./svc/forwardtests go/lint
	@$(MAKE) -C ./svc/indicators go/lint
	@$(MAKE) -C ./svc/ticks go/lint

.PHONY: lint
lint: go/lint ## Lint the golang code

.PHONY: go/test/unit
go/test/unit: ## Launch golang unit tests
	@$(MAKE) -C ./clients/go go/test/unit
	@$(MAKE) -C ./cmd/cryptellation go/test/unit
	@$(MAKE) -C ./cmd/cryptellation-tui go/test/unit
	@$(MAKE) -C ./examples/go go/test/unit
	@$(MAKE) -C ./pkg go/test/unit
	@$(MAKE) -C ./svc/backtests go/test/unit
	@$(MAKE) -C ./svc/candlesticks go/test/unit
	@$(MAKE) -C ./svc/exchanges go/test/unit
	@$(MAKE) -C ./svc/forwardtests go/test/unit
	@$(MAKE) -C ./svc/indicators go/test/unit
	@$(MAKE) -C ./svc/ticks go/test/unit

.PHONY: test/unit
test/unit: go/test/unit ## Launch unit tests

.PHONY: go/test/integration
go/test/integration: ## Launch golang integration tests
	@$(MAKE) -C ./svc/backtests go/test/integration
	@$(MAKE) -C ./svc/candlesticks go/test/integration
	@$(MAKE) -C ./svc/exchanges go/test/integration
	@$(MAKE) -C ./svc/forwardtests go/test/integration
	@$(MAKE) -C ./svc/indicators go/test/integration
	@$(MAKE) -C ./svc/ticks go/test/integration

.PHONY: test/integration 
test/integration: local/up go/test/integration ## Launch integration tests

.PHONY: go/test/end-to-end
go/test/end-to-end: ## Launch golang end-to-end tests
	@$(MAKE) -C ./svc/backtests go/test/end-to-end
	@$(MAKE) -C ./svc/candlesticks go/test/end-to-end
	@$(MAKE) -C ./svc/exchanges go/test/end-to-end
	@$(MAKE) -C ./svc/forwardtests go/test/end-to-end
	@$(MAKE) -C ./svc/indicators go/test/end-to-end
	@$(MAKE) -C ./svc/ticks go/test/end-to-end

.PHONY: test/end-to-end
test/end-to-end: local/up go/test/end-to-end ## Launch end-to-end tests

.PHONY: test
test: test/unit test/integration test/end-to-end ## Launch tests
