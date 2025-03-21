.PHONY: go/generate
go/generate: ## Generate the code
	@go generate ./...

.PHONY: go/lint
go/lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run ./...

.PHONY: go/test
go/test: go/test/unit go/test/integration go/test/end-to-end ## Launch all tests

.PHONY: go/test/unit
go/test/unit: ## Launch unit tests
	@go test $$(go list ./... | grep -v -e /activities -e /test)

.PHONY: go/test/integration
go/test/integration: docker/env/up ## Launch integration tests
	@go test $$(go list ./pkg/domains/... | grep /activities)

.PHONY: go/test/end-to-end
go/test/end-to-end: docker/worker/up ## Launch end-to-end tests
	@go test ./test/go/...