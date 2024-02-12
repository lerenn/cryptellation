DOCKER_COMPOSE_CMD := docker compose -f ./deployments/docker-compose.yaml 

.PHONY: serve
serve:
	@$(DOCKER_COMPOSE_CMD) up