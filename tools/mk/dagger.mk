DAGGER_CMD         := dagger call -m ./build/ci

.PHONY: dagger/check
dagger/check: dagger/check-generation dagger/check-todos dagger/lint dagger/tests/unit ## Run all checks through Dagger

.PHONY: dagger/check-generation
dagger/check-generation: ## Run all checks for generated code through Dagger
	@$(DAGGER_CMD) check-generation --source-dir=. stdout

.PHONY: dagger/check-todos
dagger/check-todos: ## Run all checks for TODOs through Dagger
	@$(DAGGER_CMD) check-todos --source-dir=. stdout

.PHONY: dagger/create-release
dagger/create-release: ## Create a release through Dagger
	@$(DAGGER_CMD) create-release --source-dir=. --ssh-private-key-file=file://~/.ssh/id_rsa

.PHONY: dagger/develop
dagger/develop: ## Run Dagger develop on all Dagger modules
	@dagger develop -m ./pkg/dagger
	@dagger develop -m ./build/ci/dagger

.PHONY: dagger/lint
dagger/lint: ## Run all linters through Dagger
	@$(DAGGER_CMD) linter --source-dir=. stdout

.PHONY: dagger/publish-release
dagger/publish-release: ## Publish a release through Dagger
	@$(DAGGER_CMD) publish-release --source-dir=. --ssh-private-key-file=file://~/.ssh/id_rsa

.PHONY: dagger/tests/unit
dagger/tests/unit: ## Run all unit tests through Dagger
	@$(DAGGER_CMD) unit-tests --source-dir=. stdout