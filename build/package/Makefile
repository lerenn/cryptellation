.DEFAULT_GOAL     	= help

DOCKER_IMAGE_PREFIX = digital-feather/cryptellation

RELEASE_VER			= $(shell grep "version = " ../../pkg/version/version.go | head -n 1 | cut -d \" -f 2)

# Force shell to bash for Gitlab CI
SHELL             := /bin/bash

.PHONY: build/build
build/build:
	@DOCKER_BUILDKIT=1 docker build \
		-f ./build.dockerfile \
		-t $(DOCKER_IMAGE_PREFIX):build-$(RELEASE_VER) \
		-t $(DOCKER_IMAGE_PREFIX):build \
		../..

.PHONY: build/run
build/run:
	@DOCKER_BUILDKIT=1 docker build \
		-f ./run.dockerfile \
		--build-arg BUILD_IMAGE=$(DOCKER_IMAGE_PREFIX):build-$(RELEASE_VER) \
		-t $(DOCKER_IMAGE_PREFIX):$(RELEASE_VER) \
		-t $(DOCKER_IMAGE_PREFIX):latest \
		../..

.PHONY: build 
build: build/build build/run ## Build local docker image

.PHONY: push
push: build ## Push docker image
	@docker push $(DOCKER_IMAGE):${RELEASE_VER}

.PHONY: push/local
push/local: build
	@docker tag $(DOCKER_IMAGE_PREFIX):$(RELEASE_VER) localhost:5000/$(DOCKER_IMAGE_PREFIX):$(RELEASE_VER)
	@docker push localhost:5000/$(DOCKER_IMAGE_PREFIX):$(RELEASE_VER)
	@docker tag $(DOCKER_IMAGE_PREFIX):latest localhost:5000/$(DOCKER_IMAGE_PREFIX):latest
	@docker push localhost:5000/$(DOCKER_IMAGE_PREFIX):latest

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = "\:.*?## "} /^[a-zA-Z0-9_\/-]+\:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
