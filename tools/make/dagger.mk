.PHONY: dagger/check
dagger/check: ## Run all checks through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) check \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--secrets-file=file:$(PROJECT_ROOT_PATH)/.credentials.env \
		stdout

.PHONY: dagger/lint
dagger/lint: ## Run all linters through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) lint --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/check-generation
dagger/check-generation: ## Run all checks for generated code through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) check-generation --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/unit-tests
dagger/unit-tests: ## Run all unit tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) unit-tests --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/integration-tests
dagger/integration-tests: ## Run all integration tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) integration-tests \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--secrets-file=file:$(PROJECT_ROOT_PATH)/.credentials.env \
		stdout

.PHONY: dagger/end-to-end-tests
dagger/end-to-end-tests: ## Run all end-to-end tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) end-to-end-tests \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--secrets-file=file:$(PROJECT_ROOT_PATH)/.credentials.env \
		stdout

.PHONY: dagger/publish
dagger/publish: ## Publish the project through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) publish \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--ssh-private-key-file=file:~/.ssh/id_rsa