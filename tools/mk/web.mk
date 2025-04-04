.PHONY: web/ui/build
web/ui/build: web/ui/install ## Build the ui
	@cd web/ui && yarn run build

.PHONY: web/ui/install
web/ui/install: ## Install the ui
	@cd web/ui && yarn install

.PHONY: web/ui/start
web/ui/start: web/ui/install ## Start the ui
	@cd web/ui && yarn start