DOCKER_IMAGE_PREFIX := lerenn/cryptellation

.DEFAULT_GOAL       := help

.PHONY: clean
clean: ## Clean everything
	@rm *.tgz || true

.PHONE: deploy/custom
deploy/custom: ## Deploy cryptellation helm chart on custom environment
	@helm upgrade --install cryptellation ./cryptellation -f ./custom.yml \
		-n cryptellation --create-namespace

.PHONE: deploy/local
deploy/local: ## Deploy cryptellation helm chart locally
	@helm upgrade --install cryptellation ./cryptellation -f ./local.yml

.PHONY: delete
delete: ## Delete cryptellation helm chart deployment
	@helm delete cryptellation || true


.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

.PHONY: package
package: ## Package the helm chart
	@helm package cryptellation
	
.PHONE: template/local
template/local: ## Template cryptellation helm chart locally
	@helm template cryptellation ./cryptellation -f ./local.yml --debug

