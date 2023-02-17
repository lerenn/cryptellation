.DEFAULT_GOAL := help

DOCKER_REPO    = digital-feather/cryptellation
IMAGE_NAME     = insula-api
DOCKER_IMAGE   = ${DOCKER_REPO}/$(IMAGE_NAME)

RELEASE_VER    = $(shell grep "version = " ../../pkg/version/version.go | head -n 1 | cut -d \" -f 2)

# Force shell to bash for Gitlab CI
SHELL         := /bin/bash

.PHONY: build
build: ## Build local docker image
	@DOCKER_BUILDKIT=1 docker build -f ./Dockerfile --tag $(DOCKER_IMAGE):${RELEASE_VER} ../..

.PHONY: push
push: ## Push docker image
	@docker push $(DOCKER_IMAGE):${RELEASE_VER}

.PHONY: clean
clean: ## Clean all docker images
	@docker rmi ${DOCKER_IMAGE}:${RELEASE_VER} || true

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = "\:.*?## "} /^[a-zA-Z0-9_\/-]+\:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
