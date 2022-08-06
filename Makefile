.PHONY: all
.DEFAULT_GOAL := help

docker/build: ## Build docker images
	@echo -e "\e[94m[Build docker images]\e[0m"
	@DOCKER_BUILDKIT=1 docker-compose build

docker/run: docker/build ## Run with docker
	@echo -e "\e[94m[Running locally]\e[0m"
	@docker-compose up

docker/clean: test/clean ## Clean remaining docker containers
	@echo -e "\e[94m[Cleaning remaining containers]\e[0m"
	@docker-compose down

clean: docker/clean test/clean ## Clean everything

test: test/clean test/unit test/integration test/end-to-end ## Run all tests

test/clean: ## Clean the results from tests
	@echo -e "\e[94m[Cleaning tests results and containers]\e[0m"
	@docker-compose -f tests/integration/docker-compose.yml down
	@rm -f cover.out

test/end-to-end: ## Run end-to-end tests only
	@echo -e "\e[94m[Running end-to-end tests]\e[0m"
	@echo TODO

test/integration: ## Run integration tests only
	@echo -e "\e[94m[Running integration tests]\e[0m"
	@docker-compose -f tests/integration/docker-compose.yml run tests
	@docker-compose -f tests/integration/docker-compose.yml down
	@echo -e "\e[94m[Displaying results]\e[0m"
	@go tool cover -func cover.out
	@rm -f cover.out

test/unit: ## Run unit tests only
	@echo -e "\e[94m[Running unit tests]\e[0m"
	@go test $(shell go list ./... | grep -v -e /adapters -e /service$$ ) -coverprofile cover.out
	@echo -e "\e[94m[Displaying results]\e[0m"
	@go tool cover -func cover.out
	@rm -f cover.out

proto: proto/golang proto/python ## Generate protobuf server/clients code

proto/golang:
	@echo -e "\e[94m[Generating Golang protobuf code]\e[0m"
	@./.make/proto/golang.sh backtests candlesticks exchanges livetests ticks

proto/python:
	@echo -e "\e[94m[Generating Python protobuf code]\e[0m"
	@./.make/proto/python.sh backtests candlesticks exchanges livetests ticks

lint: lint/golang ## Lint the server and clients code

lint/golang:
	@echo -e "\e[94m[Linting Golang code]\e[0m"
	@./.make/lint/golang.sh

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
