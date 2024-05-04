DAGGER_CMD        := dagger run go run ./tools/ci

ifndef TAGS
	TAGS=""
endif

.PHONY: dagger/ci
dagger/ci: ## Execute all basic CI steps
	@$(DAGGER_CMD)

.PHONY: dagger/generate
dagger/generate: ## Generate code files
	@$(DAGGER_CMD) generate

.PHONY: dagger/lint
dagger/lint: ## Lint code
	@$(DAGGER_CMD) lint

.PHONY: dagger/publish
dagger/publish: ## Publish new tag on git, docker hub, etc.
	@$(DAGGER_CMD) publish --tags ${TAGS}

.PHONY: dagger/dagger/test
dagger/test: test/unit test/integration test/end-to-end ## Launch tests 

.PHONY: dagger/test/unit
dagger/test/unit: ## Launch unit tests
	@$(DAGGER_CMD) test --type=unit

.PHONY: dagger/test/integration
dagger/test/integration: ## Launch integration tests
	@$(DAGGER_CMD) test --type=integration

.PHONY: dagger/test/end-to-end
dagger/test/end-to-end: ## Launch end-to-end tests
	@$(DAGGER_CMD) test --type=end-to-end

.PHONY: dagger/update
dagger/update: ## Update the dependencies
	@$(DAGGER_CMD) update
