.PHONY: all
.DEFAULT_GOAL := help

CLIENTS := $(shell find clients -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")
SERVICES := $(shell find services -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")

docker/build: ## Build docker images
	@DOCKER_BUILDKIT=1 docker-compose build

docker/run: docker/build ## Run with docker
	@docker-compose up

docker/clean: ## Clean remaining docker containers
	@docker-compose down

clean: docker/clean ## Clean everything
	@for CLIENT in $(CLIENTS); do $(MAKE) -C clients/$$CLIENT clean || exit $?; done
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE clean || exit $?; done

proto: ## Generate protobuf code
	@for CLIENT in $(CLIENTS); do $(MAKE) -C clients/$$CLIENT proto || exit $?; done
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE proto || exit $?; done

lint: ## Lint the code
	@for CLIENT in $(CLIENTS); do $(MAKE) -C clients/$$CLIENT lint || exit $?; done
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE lint || exit $?; done

test: ## Test services
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE test || exit $?; done

run: docker/run ## Run services (with docker)

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
