.DEFAULT_GOAL := help

MAKEFLAGS     += --no-print-directory

export NATS_HOST=localhost
export NATS_PORT=4222

.PHONY: run/local
run/local: ## Run end-to-end tests with port-forwarding from minikube
	@kubectl port-forward service/cryptellation-nats 4222:4222 &> /dev/null & FORWARD_PID="$$!"; \
		$(MAKE) run; \
		kill -TERM $$FORWARD_PID

.PHONY: run
run: ## Run end-to-end tests only
	@go test -count=1 ./... -coverprofile cover.out -v
	@go tool cover -func cover.out
	@rm cover.out

.PHONY: clean
clean: ## Clean remaining ressources
	# Nothing to do

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
