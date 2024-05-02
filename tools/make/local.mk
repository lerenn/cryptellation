DOCKER_COMPOSE_CMD := docker compose -f ./deployments/docker-compose.yaml 

.PHONY: dev/down
dev/down: ## Stop the development environment
	@$(DOCKER_COMPOSE_CMD) down

.PHONY: up
dev/up: ## Start the development environment
	@$(DOCKER_COMPOSE_CMD) up \
		--attach backtests-api \
		--attach candlesticks-api \
		--attach exchanges-api \
		--attach indicators-api \
		--attach ticks-api
