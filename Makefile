.PHONY: all
.DEFAULT_GOAL := help

SHELL=bash

CLIENTS := $(shell find clients -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")
SERVICES := $(shell find services -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")

SHORT_COMMIT_SHA := $(shell git rev-parse --short HEAD)

DOCKER_COMPOSE := docker-compose -p cryptellation $(foreach var,$(SERVICES),-f services/$(var)/docker-compose.yml)
DOCKER_IMAGE_TAG=$(SHORT_COMMIT_SHA)
DOCKER_BUILDKIT=1

export

docker/clean: ## Clean remaining docker containers
	$(DOCKER_COMPOSE) down

docker/build: ## Build docker images
	$(DOCKER_COMPOSE) build

docker/push: docker/build ## Push docker images
	@git diff-index --quiet HEAD || (echo "ERROR: Some files have been modified. Please commit before pushing."; exit 1)
	$(DOCKER_COMPOSE) push

docker/run: docker/build ## Run with docker
	$(DOCKER_COMPOSE) up

docker/status: ## Display docker status
	$(DOCKER_COMPOSE) ps

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
