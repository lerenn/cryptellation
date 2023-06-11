.PHONY: all
.DEFAULT_GOAL  := help

NAMESPACE      := cryptellation

DOCKER_COMPOSE := docker-compose -p $(NAMESPACE)

.PHONY: clean
clean: ## Clean up tests
	@$(DOCKER_COMPOSE) -f ./docker-compose.yml down
	@docker rmi $(NAMESPACE)_tests || true
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
