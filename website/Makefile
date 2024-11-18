.PHONY: all
.DEFAULT_GOAL := help

.PHONY: serve
serve: ## Serve the documentation website locally
	@go run github.com/gohugoio/hugo@v0.117.0 server

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
