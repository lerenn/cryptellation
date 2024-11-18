DOCKER_IMAGE        := lerenn/cryptellation

GIT_COMMIT_SHA      := $(shell git rev-parse HEAD)
GIT_LAST_BRANCH_TAG := $(shell git describe --abbrev=0 --tags)

.PHONY: docker/build
docker/build: ## Build the docker image
	@docker buildx create --use --name=crossplatform --node=crossplatform && \
	docker buildx build \
		--file ./Dockerfile \
		--output "type=docker,push=false" \
		--tag $(DOCKER_IMAGE):devel \
		../..

.PHONY: docker/publish
docker/publish: ## Publish the docker image
	@docker buildx create --use --name=crossplatform --node=crossplatform && \
	docker buildx build \
		--file ./build/package/Dockerfile \
		--platform linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64,linux/ppc64le,linux/s390x \
		--output "type=image,push=true" \
		--tag $(DOCKER_IMAGE):$(GIT_COMMIT_SHA) \
		--tag $(DOCKER_IMAGE):$(GIT_LAST_BRANCH_TAG) \
		--tag $(DOCKER_IMAGE):latest \
		../..