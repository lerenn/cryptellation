K8S_VERSION         := 1.26.1
MINIKUBE_NODES      := 3
NODE_MEMORY         := 2g
CLUSTER_NAME        := cryptellation-local

.PHONY: minikube/up
minikube/up: ## Deploy minikube cluster
	@minikube start \
		--cpus=2 \
		--memory=${NODE_MEMORY} \
		--nodes=${MINIKUBE_NODES} \
		--profile ${CLUSTER_NAME} \
		--kubernetes-version=v${K8S_VERSION} \
		--extra-config=kubelet.cluster-dns=10.96.0.10
	@minikube -p $(CLUSTER_NAME) addons enable registry
	@helm repo add otel-lgtm-dev https://lerenn.github.io/packages/helm/otel-lgtm-dev
	@helm upgrade --install lgtm otel-lgtm-dev/otel-lgtm-dev \
		-n telemetry --create-namespace

.PHONY: minikube/status
minikube/status: ## Checks the minikube status
	@minikube status --profile ${CLUSTER_NAME}

.PHONY: minikube/stop
minikube/stop: ## Stop the current minikube cluster
	@minikube stop -p ${CLUSTER_NAME}

.PHONY: minikube/destroy
minikube/destroy: ## Destroy the minikube cluster
	@minikube delete --profile ${CLUSTER_NAME}

.PHONY: minikube/expose
minikube/expose: ## Expose the ports to 8080
	@kubectl port-forward -n telemetry service/lgtm-grafana 8080:80 & \
		kubectl port-forward service/cryptellation-nats 4222:4222 & \
		kubectl port-forward --namespace kube-system service/registry 5000:80 & \
		wait