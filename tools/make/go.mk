.PHONY: go/generate
go/generate: ## Generate the code
	@go generate ./...

.PHONY: go/lint
go/lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55 run ./...

.PHONY: go/test/unit
go/test/unit: ## Launch unit tests	
	@go test $$(go list ./... | grep -v -e /adapters -e /test)

.PHONY: go/test/integration
go/test/integration: ## Launch integration tests
	@go run ./cmd/data migrations migrate
	@go test ./internal/adapters/...

.PHONY: go/test/end-to-end
go/test/end-to-end: ## Launch end-to-end tests
	@go run ./cmd/data migrations migrate
	@go test ./test/...