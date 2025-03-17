.PHONY: code/check
code/check: generate lint test ## Generate, lint and test the code

.PHONY: code/check-todos
code/check-todos: ## Check the todos in the code
	@go run ./tools/invtodos .

.PHONY: code/generate
code/generate: ## Generate the code
	@go generate ./...

.PHONY: code/lint
code/lint: ## Lint the code
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 run ./...

.PHONY: code/test
code/test: code/test/unit code/test/integration code/test/end-to-end ## Launch all tests

.PHONY: code/test/unit
code/test/unit: ## Launch unit tests
	@go test $$(go list ./... | grep -v -e /activities -e /test)

.PHONY: code/test/integration
code/test/integration: docker/env/up ## Launch integration tests
	@go test $$(go list ./pkg/domains/... | grep /activities)

.PHONY: code/test/end-to-end
code/test/end-to-end: docker/worker/up ## Launch end-to-end tests
	@go test ./test/...