DOCKER_COMPOSE_CMD := docker compose -f ./deployments/docker-compose.yaml 

.PHONY: build
build: docker/build ## Build the project

.PHONY: clean
clean: local/down ## Clean the project

.PHONY: generate
generate: go/generate ## Generate the code

.PHONY: local/down
local/down: ## Stop the local environment
	@$(DOCKER_COMPOSE_CMD) down

.PHONY: local/up
local/up: ## Start the local environment
	@$(DOCKER_COMPOSE_CMD) up -d

.PHONY: lint
lint: go/lint ## Lint the golang code

.PHONY: test/unit
test/unit: go/test/unit ## Launch unit tests

.PHONY: test/integration 
test/integration: local/up go/test/integration ## Launch integration tests

.PHONY: test/end-to-end
test/end-to-end: local/up go/test/end-to-end ## Launch end-to-end tests

.PHONY: test
test: test/unit test/integration ## Launch tests