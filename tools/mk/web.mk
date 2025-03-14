.PHONY: web/ui/build
web/ui/build: web/ui/install ## Build the ui
	@cd web/ui && npm run build

.PHONY: web/ui/install
web/ui/install: ## Install the ui
	@cd web/ui && npm install

.PHONY: web/ui/start
web/ui/start: web/ui/install ## Start the ui
	@cd web/ui && npm start