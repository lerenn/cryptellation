.PHONY: web/dashboard/start
web/dashboard/start: ## Start the dashboard
	@cd web/dashboard && npm start

.PHONY: web/dashboard/build
web/dashboard/build: ## Build the dashboard
	@cd web/dashboard && npm run build