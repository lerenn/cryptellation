DOCKER_IMAGE_PREFIX := lerenn/cryptellation

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

.PHONY: clean
clean: delete/local ## Clean everything

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

