KIND_CMD     := go run sigs.k8s.io/kind@v0.23.0
CLUSTER_NAME := cryptellation-cluster

.PHONY: kind/up
kind/up: ## Deploy kind cluster
	@${KIND_CMD} create cluster --config ./deployments/kind.yaml --name ${CLUSTER_NAME}

.PHONY: kind/telemetry/up
kind/telemetry/up: ## Deploy telemetry on kind cluster
	@helm repo add otel-lgtm-dev https://lerenn.github.io/packages/helm/otel-lgtm-dev
	@helm upgrade --install lgtm otel-lgtm-dev/otel-lgtm-dev \
		-n telemetry --create-namespace

.PHONY: kind/telemetry/forward
kind/telemetry/forward: ## Forward telemetry on kind cluster
	@kubectl port-forward svc/lgtm-otel-collector 4317:4317 -n telemetry

.PHONY: kind/telemetry/down
kind/telemetry/down: ## Destroy telemetry on kind cluster
	@helm uninstall lgtm -n telemetry
	@kubectl delete ns telemetry

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