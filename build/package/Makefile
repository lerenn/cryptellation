.DEFAULT_GOAL     := help

DOCKER_IMAGE_PREFIX = digital-feather/cryptellation

SERVICES          := $(shell find ../../internal -mindepth 1 -maxdepth 1 -type d | xargs -I{} basename "{}")

# Force shell to bash for Gitlab CI
SHELL             := /bin/bash

.PHONY: build
build: ## Build local docker image
	for SERVICE in $(SERVICES); do \
		VERSION=$$(grep "SemVer: " ../../internal/$${SERVICE}/version.go | head -n 1 | cut -d \" -f 2); \
		DOCKER_BUILDKIT=1 docker build -f ./Dockerfile --build-arg COMMAND=cryptellation-$${SERVICE} --tag $(DOCKER_IMAGE_PREFIX)-$${SERVICE}:$${VERSION} ../..; \
		done

.PHONY: push
push: ## Push docker image
	@docker push $(DOCKER_IMAGE):${RELEASE_VER}

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = "\:.*?## "} /^[a-zA-Z0-9_\/-]+\:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
