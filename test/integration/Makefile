.PHONY: all
.DEFAULT_GOAL := help

DOCKER_COMPOSE := docker-compose -p cryptellation

.PHONY: clean
clean: ## Clean up tests
	@$(DOCKER_COMPOSE) -f ./docker-compose.yml down
	@rm -f cover.out

.PHONY: run
run: ## Perform integration tests
	@$(DOCKER_COMPOSE) -f ./docker-compose.yml build
	@$(DOCKER_COMPOSE) -f ./docker-compose.yml run tests

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
