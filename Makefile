.PHONY: all
.DEFAULT_GOAL := help

GOLANG_SERVICES := backtests candlesticks exchanges livetests ticks
SERVICES := $(GOLANG_SERVICES)

docker/build: ## Build docker images
	@echo -e "\e[94m[Build docker images]\e[0m"
	@DOCKER_BUILDKIT=1 docker-compose build

docker/run: docker/build ## Run with docker
	@echo -e "\e[94m[Running locally]\e[0m"
	@docker-compose up

docker/clean: ## Clean remaining docker containers
	@echo -e "\e[94m[Cleaning remaining containers]\e[0m"
	@docker-compose down

clean: docker/clean ## Clean everything

proto: proto/golang proto/python ## Generate protobuf server/clients code

proto/golang:
	@echo -e "\e[94m[Generating Golang protobuf code]\e[0m"
	@./.make/proto/golang.sh $(SERVICES)

proto/python:
	@echo -e "\e[94m[Generating Python protobuf code]\e[0m"
	@./.make/proto/python.sh $(SERVICES)

lint: lint/golang ## Lint the server and clients code

lint/golang:
	@./.make/lint/golang.sh

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
