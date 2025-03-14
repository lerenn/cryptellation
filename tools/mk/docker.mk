DOCKER_CMD          := docker
DOCKER_COMPOSE_CMD  := $(DOCKER_CMD) compose -f ./deployments/docker-compose/docker-compose.yaml
DOCKER_IMAGE        := lerenn/cryptellation
DOCKER_CFG_PATH     := ./build/package

GIT_COMMIT_SHA      := $(shell git rev-parse HEAD)
GIT_LAST_BRANCH_TAG := $(shell git describe --abbrev=0 --tags)

.PHONY: docker/build
docker/build: ## Build the docker image
	@$(DOCKER_CMD) buildx create --use --name=cryptellation --node=cryptellation
	@$(DOCKER_CMD) buildx build \
		--file $(DOCKER_CFG_PATH)/Dockerfile \
		--output "type=docker,push=false" \
		--tag $(DOCKER_IMAGE):devel \
		.

.PHONY: docker/clean
docker/clean: docker/env/down ## Clean the docker environment
	@$(DOCKER_CMD) rmi $(DOCKER_IMAGE):devel || true
	@$(DOCKER_CMD) buildx rm cryptellation || true
	@$(DOCKER_COMPOSE_CMD) --profile worker --profile migrator down --volumes --remove-orphans

.PHONY: docker/env/down 
docker/env/down: ## Stop the dependencies in local environment
	@$(DOCKER_COMPOSE_CMD) down

.PHONY: docker/env/up
docker/env/up: ## Start the dependencies in local environment
	@$(DOCKER_COMPOSE_CMD) up -d
	@$(DOCKER_COMPOSE_CMD) run migrator

.PHONY: docker/publish
docker/publish: ## Publish the docker image
	@$(DOCKER_CMD) buildx create --use --name=cryptellation --node=cryptellation
	@$(DOCKER_CMD) buildx build \
		--file $(DOCKER_CFG_PATH)/Dockerfile \
		--platform linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64,linux/ppc64le,linux/s390x \
		--output "type=image,push=true" \
		--tag $(DOCKER_IMAGE):$(GIT_COMMIT_SHA) \
		--tag $(DOCKER_IMAGE):$(GIT_LAST_BRANCH_TAG) \
		--tag $(DOCKER_IMAGE):latest \
		.

.PHONY: docker/worker/down
docker/worker/down: ## Start a cryptellation worker in local environment
	@$(DOCKER_COMPOSE_CMD) --profile worker down

.PHONY: docker/worker/up
docker/worker/up: docker/env/up ## Start a cryptellation worker in local environment
	@$(DOCKER_COMPOSE_CMD) --profile worker up -d
