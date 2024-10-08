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

.PHONY: dagger/tests/unit
dagger/tests/unit: ## Run all unit tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) unit-tests --source-dir=$(PROJECT_ROOT_PATH) stdout

.PHONY: dagger/tests/integration
dagger/tests/integration: ## Run all integration tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) integration-tests \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--secrets-file=file:$(PROJECT_ROOT_PATH)/.credentials.env \
		stdout

.PHONY: dagger/tests/end-to-end
dagger/tests/end-to-end: ## Run all end-to-end tests through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) end-to-end-tests \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--secrets-file=file:$(PROJECT_ROOT_PATH)/.credentials.env \
		stdout

.PHONY: dagger/tests
dagger/tests: dagger/tests/unit dagger/tests/integration dagger/tests/end-to-end ## Run all tests through Dagger

.PHONY: dagger/create-release
dagger/create-release: ## Create a release of the new version through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) create-release \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--ssh-private-key-file=file:~/.ssh/id_rsa

.PHONY: dagger/publish-release
dagger/publish-release: ## Publish the new release version through Dagger
	@dagger call -m $(DAGGER_MODULE_PATH) publish-release \
		--source-dir=$(PROJECT_ROOT_PATH) \
		--ssh-private-key-file=file:~/.ssh/id_rsa
