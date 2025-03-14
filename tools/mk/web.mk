.PHONY: web/ui/start
web/ui/start: ## Start the ui
	@cd web/ui && npm start

.PHONY: web/ui/build
web/ui/build: ## Build the ui
	@cd web/ui && npm run build