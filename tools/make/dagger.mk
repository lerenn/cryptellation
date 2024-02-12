DAGGER_CMD        := dagger run go run ./tools/ci

.PHONY: ci
ci: ## Execute all basic CI steps
	@$(DAGGER_CMD)

.PHONY: generate
generate: ## Generate code files
	@$(DAGGER_CMD) generate

.PHONY: lint
lint: ## Lint code
	@$(DAGGER_CMD) lint

.PHONY: publish
publish: ## Publish docker images publically
	@$(DAGGER_CMD) publish

.PHONY: test
test: test/unit test/integration test/end-to-end ## Launch tests 

.PHONY: test/unit
test/unit: ## Launch unit tests
	@$(DAGGER_CMD) test --type=unit

.PHONY: test/integration
test/integration: ## Launch integration tests
	@$(DAGGER_CMD) test --type=integration

.PHONY: test/end-to-end
test/end-to-end: ## Launch end-to-end tests
	@$(DAGGER_CMD) test --type=end-to-end

.PHONY: update
update: ## Update the dependencies
	@$(DAGGER_CMD) update
