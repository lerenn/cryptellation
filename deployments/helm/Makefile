DOCKER_IMAGE_PREFIX := lerenn/cryptellation

.PHONY: clean
clean: ## Clean everything
	@rm *.tgz || true

.PHONY: deploy/local
deploy/local: ## Deploy cryptellation on local environment
	@kubectl port-forward --namespace kube-system service/registry 5000:80 &> /dev/null & FORWARD_PID="$$!"; \
		$(MAKE) -C ../../build/package push/local; \
		kill -TERM $$FORWARD_PID
	@make deploy

.PHONE: deploy
deploy: ## Deploy cryptellation helm chart
	@helm upgrade --install cryptellation ./cryptellation -f ./local.yml

.PHONY: delete
delete: ## Delete cryptellation helm chart deployment
	@helm delete cryptellation || true

.PHONY: package
package: ## Package the helm chart
	@helm package cryptellation

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

