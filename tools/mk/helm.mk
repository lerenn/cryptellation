DOCKER_IMAGE_PREFIX := lerenn/cryptellation

HELM_CMD      := helm
HELM_CFG_PATH := ./deployments/helm

.PHONY: helm/clean
helm/clean: ## Clean everything
	@rm -f $(HELM_CFG_PATH)/*.tgz

.PHONY: helm/dependencies/update
helm/dependencies/update: ## Update helm dependencies
	@$(HELM_CMD) dependency update $(HELM_CFG_PATH)/cryptellation

.PHONE: helm/deploy/custom
helm/deploy/custom: helm/dependencies/update ## Deploy cryptellation helm chart on custom environment
	@$(HELM_CMD) upgrade --install cryptellation $(HELM_CFG_PATH)/cryptellation -f $(HELM_CFG_PATH)/custom.yml \
		-n cryptellation --create-namespace

.PHONE: helm/deploy/local
helm/deploy/local: helm/dependencies/update ## Deploy cryptellation helm chart locally
	@$(HELM_CMD) upgrade --install cryptellation $(HELM_CFG_PATH)/cryptellation -f $(HELM_CFG_PATH)/local.yml

.PHONY: helm/delete
helm/delete: ## Delete cryptellation helm chart deployment
	@$(HELM_CMD) delete cryptellation || true

.PHONY: helm/package
helm/package: helm/dependencies/update ## Package the helm chart
	@$(HELM_CMD) package cryptellation
	
.PHONE: helm/template/local
helm/template/local: helm/dependencies/update ## Template cryptellation helm chart locally
	@$(HELM_CMD) template cryptellation $(HELM_CFG_PATH)/cryptellation -f $(HELM_CFG_PATH)/local.yml --debug

