.PHONY: all
.DEFAULT_GOAL := help

SHELL=bash

CLIENTS := $(shell find clients -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")
SERVICES := $(shell find services -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")

SHORT_COMMIT_SHA := $(shell git rev-parse --short HEAD)

DOCKER_IMAGE_TAG=$(SHORT_COMMIT_SHA)
DOCKER_BUILDKIT=1

docker/build: ## Build docker images
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE docker/build || exit $?; done

docker/push: docker/build ## Push docker images
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE docker/push || exit $?; done

clean:  ## Clean everything
	@for CLIENT in $(CLIENTS); do $(MAKE) -C clients/$$CLIENT clean || exit $?; done
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE clean || exit $?; done

generate: ## Generate code
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE generate || exit $?; done

lint: ## Lint the code
	@for CLIENT in $(CLIENTS); do $(MAKE) -C clients/$$CLIENT lint || exit $?; done
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE lint || exit $?; done

test: ## Test services
	@for SERVICE in $(SERVICES); do $(MAKE) -C services/$$SERVICE test || exit $?; done

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
