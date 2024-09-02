KIND_CMD     := go run sigs.k8s.io/kind@v0.23.0
HELM_CMD     := helm
KUBECTL_CMD  := kubectl
CLUSTER_NAME := cryptellation-cluster

.PHONY: kind/up
kind/up: ## Deploy kind cluster
	@${KIND_CMD} create cluster --config ./deployments/kind.yaml --name ${CLUSTER_NAME}

.PHONY: kind/telemetry/up
kind/telemetry/up: ## Deploy telemetry on kind cluster
	@$(HELM_CMD) repo update
	@$(HELM_CMD) repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	@$(HELM_CMD) repo add uptrace https://charts.uptrace.dev
	@$(HELM_CMD) upgrade --install otel-collector open-telemetry/opentelemetry-collector \
		-f ./tools/uptrace/kind/otel-collector.yaml \
		-n telemetry --create-namespace
	@$(HELM_CMD) upgrade --install uptrace uptrace/uptrace \
		-f ./tools/uptrace/kind/uptrace.yaml \
		-n telemetry --create-namespace

.PHONY: kind/telemetry/forward
kind/telemetry/forward: ## Forward telemetry on kind cluster
	@$(KUBECTL_CMD) port-forward svc/uptrace 14318:14318 -n telemetry

.PHONY: kind/telemetry/down
kind/telemetry/down: ## Destroy telemetry on kind cluster
	@$(HELM_CMD) uninstall uptrace -n telemetry || true
	@$(HELM_CMD) uninstall otel-collector -n telemetry || true
	@$(KUBECTL_CMD) delete ns telemetry || true

.PHONY: kind/cryptellation/load-images
kind/cryptellation/load-images: docker/build ## Load images into kind cluster
	@${KIND_CMD} load docker-image --name ${CLUSTER_NAME} \
		lerenn/cryptellation-backtests:devel \
		lerenn/cryptellation-candlesticks:devel \
		lerenn/cryptellation-exchanges:devel \
		lerenn/cryptellation-forwardtests:devel \
		lerenn/cryptellation-indicators:devel \
		lerenn/cryptellation-ticks:devel

.PHONY: kind/cryptellation/deploy
kind/cryptellation/deploy: kind/cryptellation/load-images ## Deploy cryptellation on kind cluster
	@$(MAKE) -C deployments/helm deploy/local

.PHONY: kind/cryptellation/forward
kind/cryptellation/forward: ## Forward cryptellation on kind cluster
	@kubectl port-forward service/cryptellation-nats 4222:4222

.PHONY: kind/cryptellation/delete
kind/cryptellation/delete: ## Delete cryptellation on kind cluster
	@$(MAKE) -C deployments/helm delete

.PHONY: kind/down
kind/down: ## Destroy kind cluster
	@${KIND_CMD} delete cluster --name ${CLUSTER_NAME}