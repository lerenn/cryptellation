.DEFAULT_GOAL := help

MAKEFLAGS     += --no-print-directory

export NATS_HOST=localhost
export NATS_PORT=4222

.PHONY: run
run: ## Run end-to-end tests with port-forwarding from minikube
	@kubectl port-forward --namespace event-broker service/nats 4222:4222 &> /dev/null & FORWARD_PID="$$!"; \
		$(MAKE) test; \
		kill -TERM $$FORWARD_PID

.PHONY: clean
clean: ## Clean remaining ressources
	# Nothing to do

.PHONY: test
test: ## Run end-to-end tests only
	@go test -count=1 ./...

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
