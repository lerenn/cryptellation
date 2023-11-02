K8S_VERSION         := 1.26.1
MINIKUBE_NODES      := 3
NODE_MEMORY         := 2g
CLUSTER_NAME        := cryptellation-devel

.PHONY: minikube/deploy
minikube/deploy: ## Deploy minikube cluster
	@minikube start \
		--cpus=2 \
		--memory=${NODE_MEMORY} \
		--nodes=${MINIKUBE_NODES} \
		--profile ${CLUSTER_NAME} \
		--kubernetes-version=v${K8S_VERSION} \
		--extra-config=kubelet.cluster-dns=10.96.0.10
	@minikube -p $(CLUSTER_NAME) addons enable registry

.PHONY: minikube/status
minikube/status: ## Checks the minikube status
	@minikube status --profile ${CLUSTER_NAME}

.PHONY: minikube/stop
minikube/stop: ## Stop the current minikube cluster
	@minikube stop -p ${CLUSTER_NAME}

.PHONY: minikube/delete
minikube/delete: ## Delete the minikube cluster
	@minikube delete --profile ${CLUSTER_NAME}

.PHONY: minikube/deploy-lgtm
minikube/deploy-lgtm: ## Deploy a Helm chart for telemetry
	@helm repo add otel-lgtm-dev https://lerenn.github.io/packages/helm/otel-lgtm-dev
	@helm upgrade --install lgtm otel-lgtm-dev/otel-lgtm-dev \
		-n telemetry --create-namespace

.PHONY: minikube/expose-grafana
minikube/expose-grafana: ## Expose the Grafana port to 8080
	@kubectl port-forward -n telemetry service/lgtm-grafana 8080:80

.PHONY: minikube/delete-lgtm
minikube/delete-lgtm: ## Delete telemetry resources
	@helm delete -n telemetry lgtm || true
	@kubectl delete all -n telemetry --all || true
	@kubectl delete pvc -n telemetry --all || true

.PHONY: clean
clean: ## Clean everything
	@$(MAKE) minikube/delete

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'

