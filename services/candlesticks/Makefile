.PHONY: all
.DEFAULT_GOAL := help

clean: test/clean ## Clean everything

proto: ## Generate protobuf code
	@mkdir -p pkg/client/proto
	@for PROTOFILE in $(shell find . -name *.proto); do \
		echo "Generating $$PROTOFILE"; \
		protoc \
			--proto_path=./api \
			--go_out=./pkg/client/proto \
			--go_opt=paths=source_relative \
			--go-grpc_opt=require_unimplemented_servers=false \
			--go-grpc_out=./pkg/client/proto \
			--go-grpc_opt=paths=source_relative \
			"$$PROTOFILE" || exit $?; \
	done

lint: ## Lint the code
ifeq ($(shell which golangci-lint &> /dev/null; echo $$?),1)
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif
	@LOG_LEVEL=error golangci-lint run

test/build:
	@docker-compose build

test/clean:
	@docker-compose down
	@rm -f cover.out

test/run:
	@docker-compose run tests

test: test/run ## Run tests

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
