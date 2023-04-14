.PHONY: all
.DEFAULT_GOAL := help

DOCKER_COMPOSE := docker-compose -p cryptellation


.PHONY: all
all: unit integration end-to-end ## Run tests

.PHONY: clean
clean:
	$(MAKE) -C end-to-end clean
	@$(DOCKER_COMPOSE) -f ./integration/docker-compose.yml down
	@rm -f cover.out

.PHONY: end-to-end
end-to-end: ## Perform end-to-end tests
	@$(MAKE) -C end-to-end run

.PHONY: integration
integration: ## Perform integration tests
	@$(DOCKER_COMPOSE) -f ./integration/docker-compose.yml build
	@$(DOCKER_COMPOSE) -f ./integration/docker-compose.yml run tests

.PHONY: unit
unit: ## Perform unit tests
	@go test $(shell go list ../cmd/... ../pkg/... ../services/... | grep -v -e /io/) -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
